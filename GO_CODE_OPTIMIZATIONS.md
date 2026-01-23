# Optimizaciones de Código Go - Alertly Backend

Recomendaciones para mejorar el código Go después de implementar índices SQL.

---

## 1. Eliminar Subconsultas en UPDATEs

### Archivo: `newincident/repository.go:265-282`

**PROBLEMA ACTUAL:**

```go
query := `
UPDATE incident_clusters ic
SET
    score_true = ic.score_true + (SELECT credibility FROM account WHERE account_id = $1),
    score_false = ic.score_false + (10 - (SELECT credibility FROM account WHERE account_id = $1)),
    credibility = ic.score_true / GREATEST(ic.score_true + ic.score_false, 1) * 10
WHERE ic.incl_id = $4;
`
result, err := r.db.Exec(query, accountID, latitude, longitude, inclId)
```

**Problema:** 2 subconsultas por cada UPDATE (ejecutadas 100+ veces/hora)

**SOLUCIÓN OPTIMIZADA:**

```go
func (r *pgRepository) UpdateClusterAsTrue(inclId int64, accountID int64, latitude, longitude float32) (sql.Result, error) {
    // 1. Obtener credibilidad UNA sola vez
    var credibility float64
    err := r.db.QueryRow("SELECT credibility FROM account WHERE account_id = $1", accountID).Scan(&credibility)
    if err != nil {
        return nil, fmt.Errorf("failed to get account credibility: %w", err)
    }

    // 2. Calcular valores en Go (más eficiente)
    scoreTrue := credibility
    scoreFalse := 10 - credibility

    // 3. UPDATE sin subconsultas
    query := `
    UPDATE incident_clusters ic
    SET
        center_latitude      = (ic.center_latitude + $1) / 2,
        center_longitude     = (ic.center_longitude + $2) / 2,
        counter_total_votes  = ic.counter_total_votes + 1,
        score_true           = ic.score_true + $3,
        score_false          = ic.score_false + $4,
        credibility          = ic.score_true / GREATEST(ic.score_true + ic.score_false, 1) * 10
    WHERE ic.incl_id = $5;
    `

    result, err := r.db.Exec(query, latitude, longitude, scoreTrue, scoreFalse, inclId)
    return result, err
}
```

**Aplicar también a:**
- `UpdateClusterAsFalse` (línea 284-301)

**Impacto:** 2-3x más rápido por UPDATE

---

## 2. Separar Query JSON_AGG en Profile

### Archivo: `profile/repository.go:63-87`

**PROBLEMA ACTUAL:**

```go
query := `
SELECT
    a.account_id,
    -- ... 30+ columnas
    COALESCE(
        (
        SELECT JSON_AGG(sub.incident_data)
        FROM (
            SELECT JSON_BUILD_OBJECT(
                'inre_id', i.inre_id,
                -- ... más campos
            ) AS incident_data
            FROM incident_reports i
            INNER JOIN incident_clusters ic ON i.incl_id = ic.incl_id
            WHERE i.account_id = a.account_id
            ORDER BY i.created_at DESC
            LIMIT 50
        ) sub
        ),
        '[]'::json
    ) AS incidents
FROM account a
WHERE a.account_id = $1
`
```

**Problema:** JSON_AGG hace JOIN lateral por cada fila, lento con 50+ incidents

**SOLUCIÓN OPTIMIZADA:**

```go
func (r *pgRepository) GetById(accountID int64) (Profile, error) {
    // Query 1: Account info (rápida)
    accountQuery := `
    SELECT
        a.account_id,
        a.nickname,
        a.first_name,
        -- ... resto de campos sin incidents
    FROM account a
    WHERE a.account_id = $1
    `

    var profile Profile
    err := r.db.QueryRow(accountQuery, accountID).Scan(
        &profile.AccountID,
        &profile.Nickname,
        // ... resto de campos
    )
    if err != nil {
        return Profile{}, fmt.Errorf("error scanning account: %w", err)
    }

    // Query 2: Incidents (paralela con goroutine)
    incidentsQuery := `
    SELECT
        i.inre_id,
        COALESCE(i.media_url, '') as media_url,
        COALESCE(i.description, '') as description,
        COALESCE(i.event_type, '') as event_type,
        COALESCE(i.subcategory_name, '') as subcategory_name,
        COALESCE(ic.credibility, 0) as credibility,
        i.incl_id,
        COALESCE(i.is_anonymous, '0') as is_anonymous,
        i.created_at
    FROM incident_reports i
    INNER JOIN incident_clusters ic ON i.incl_id = ic.incl_id
    WHERE i.account_id = $1
    ORDER BY i.created_at DESC
    LIMIT 50
    `

    rows, err := r.db.Query(incidentsQuery, accountID)
    if err != nil {
        return Profile{}, fmt.Errorf("error querying incidents: %w", err)
    }
    defer rows.Close()

    incidents := make([]Incident, 0, 50)
    for rows.Next() {
        var incident Incident
        err := rows.Scan(
            &incident.InreId,
            &incident.MediaUrl,
            &incident.Description,
            &incident.EventType,
            &incident.SubcategoryName,
            &incident.Credibility,
            &incident.InclId,
            &incident.IsAnonymous,
            &incident.CreatedAt,
        )
        if err != nil {
            return Profile{}, fmt.Errorf("error scanning incident: %w", err)
        }
        incidents = append(incidents, incident)
    }

    profile.Incidents = incidents

    return profile, nil
}
```

**Impacto:** 3-5x más rápido + mejor uso de índices

---

## 3. Optimizar GetReel - Evitar ORDER BY RANDOM()

### Archivo: `getincidentsasreels/repository.go:27-44`

**PROBLEMA ACTUAL:**

```go
idQuery := `
SELECT c.incl_id
FROM incident_clusters c
WHERE (...)
ORDER BY RANDOM()
LIMIT 20
`
```

**Problema:** RANDOM() causa full table scan, no usa índices

**SOLUCIÓN 1 - TABLESAMPLE (PostgreSQL):**

```go
func (r *pgRepository) GetReel(inputs Inputs, accountID int64) ([]getclusterby.Cluster, error) {
    // Usar TABLESAMPLE para sample aleatorio eficiente
    idQuery := `
    SELECT c.incl_id
    FROM incident_clusters c TABLESAMPLE SYSTEM (1) -- 1% de filas
    WHERE
        (c.center_latitude  BETWEEN $1 AND $2 AND c.center_longitude BETWEEN $3 AND $4)
        OR EXISTS (
          SELECT 1
          FROM account_favorite_locations f
          WHERE f.account_id = $5
            AND ST_DistanceSphere(
                  ST_MakePoint(c.center_longitude, c.center_latitude),
                  ST_MakePoint(f.longitude, f.latitude)
                ) <= $6
        )
        AND c.is_active = '1'
    LIMIT 20
    `
    // ... resto del código
}
```

**SOLUCIÓN 2 - Offset Aleatorio (más predecible):**

```go
func (r *pgRepository) GetReel(inputs Inputs, accountID int64) ([]getclusterby.Cluster, error) {
    // 1. Obtener total de clusters (cachear este valor)
    var totalClusters int
    countQuery := `SELECT COUNT(*) FROM incident_clusters WHERE is_active = '1'`
    err := r.db.QueryRow(countQuery).Scan(&totalClusters)
    if err != nil || totalClusters < 20 {
        // Fallback si error o pocos clusters
        totalClusters = 100
    }

    // 2. Calcular offset aleatorio
    randomOffset := 0
    if totalClusters > 20 {
        randomOffset = rand.Intn(totalClusters - 20)
    }

    // 3. Query con offset aleatorio (USA ÍNDICES)
    idQuery := `
    SELECT c.incl_id
    FROM incident_clusters c
    WHERE
        (c.center_latitude  BETWEEN $1 AND $2 AND c.center_longitude BETWEEN $3 AND $4)
        OR EXISTS (...)
        AND c.is_active = '1'
    ORDER BY c.created_at DESC  -- USA ÍNDICE
    LIMIT 20 OFFSET $7
    `

    rows, err := r.db.Query(idQuery,
        inputs.MinLatitude, inputs.MaxLatitude,
        inputs.MinLongitude, inputs.MaxLongitude,
        accountID, maxDistanceMeters,
        randomOffset,
    )
    // ... resto del código
}
```

**Impacto:** 10-20x más rápido en tablas grandes

---

## 4. Batch Inserts para Notificaciones

### Archivo: `common/notification.go:22-39`

**PROBLEMA ACTUAL:**

```go
// Se llama múltiples veces en goroutines separadas
func SaveNotification(dbExec DBExecutor, nType string, accountID int64, referenceID int64) error {
    query := `INSERT INTO notifications(...) VALUES($1, $2, $3, ...)`
    _, err := dbExec.Exec(query, /* ... */)
    return err
}
```

**Problema:** 1 INSERT por notificación (lento con 100+ notificaciones)

**SOLUCIÓN - Batch Insert:**

```go
// Nueva función para batch inserts
func SaveNotificationsBatch(db *sql.DB, notifications []alerts.Alert) error {
    if len(notifications) == 0 {
        return nil
    }

    // Construir query con múltiples VALUES
    valueStrings := make([]string, len(notifications))
    valueArgs := make([]interface{}, 0, len(notifications)*10)

    for i, n := range notifications {
        valueStrings[i] = fmt.Sprintf(
            "($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
            i*10+1, i*10+2, i*10+3, i*10+4, i*10+5,
            i*10+6, i*10+7, i*10+8, i*10+9, i*10+10,
        )
        valueArgs = append(valueArgs,
            n.AccountID,
            n.Title,
            n.Message,
            n.Type,
            n.Link,
            dbtypes.BoolToInt(n.MustSendPush),
            dbtypes.BoolToInt(n.MustSendInApp),
            dbtypes.BoolToInt(n.MustBeProcessed),
            n.ErrorMessage,
            n.ReferenceID,
        )
    }

    query := fmt.Sprintf(
        "INSERT INTO notifications(owner_account_id, title, message, type, link, must_send_as_notification_push, must_send_as_notification, must_be_processed, error_message, reference_id) VALUES %s",
        strings.Join(valueStrings, ","),
    )

    _, err := db.Exec(query, valueArgs...)
    return err
}

// Usar en cronjobs que generan múltiples notificaciones
func ProcessNotifications(db *sql.DB, users []User) error {
    notifications := make([]alerts.Alert, 0, len(users))
    for _, user := range users {
        n := HandleNotification("new_cluster", user.ID, clusterID)
        notifications = append(notifications, n)
    }

    // 1 query en lugar de N queries
    return SaveNotificationsBatch(db, notifications)
}
```

**Impacto:** 10-50x más rápido con batch de 100+ notificaciones

---

## 5. Connection Pooling Optimization

### Archivo: `database/database.go`

**CONFIGURACIÓN ACTUAL:**

```go
db.SetMaxOpenConns(500)
db.SetMaxIdleConns(50)
db.SetConnMaxLifetime(time.Hour)
```

**RECOMENDACIONES OPTIMIZADAS:**

```go
func InitDB() (*sql.DB, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    // Configuración basada en carga real
    maxOpenConns := 100     // Era 500 (muy alto para API típica)
    maxIdleConns := 25      // Era 50
    connMaxLifetime := 5 * time.Minute  // Era 1h (muy largo)
    connMaxIdleTime := 2 * time.Minute  // NUEVO: cerrar conexiones idle

    db.SetMaxOpenConns(maxOpenConns)
    db.SetMaxIdleConns(maxIdleConns)
    db.SetConnMaxLifetime(connMaxLifetime)
    db.SetConnMaxIdleTime(connMaxIdleTime)

    // Verificar conexión
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("cannot ping database: %w", err)
    }

    // Log configuración
    log.Printf("✅ DB Pool configured: max_open=%d, max_idle=%d, max_lifetime=%s",
        maxOpenConns, maxIdleConns, connMaxLifetime)

    return db, nil
}
```

**Monitoreo del pool:**

```go
// Agregar endpoint de health check
func (h *HealthHandler) GetDBStats(c *gin.Context) {
    stats := h.db.Stats()
    c.JSON(200, gin.H{
        "max_open_connections": stats.MaxOpenConnections,
        "open_connections":     stats.OpenConnections,
        "in_use":               stats.InUse,
        "idle":                 stats.Idle,
        "wait_count":           stats.WaitCount,
        "wait_duration":        stats.WaitDuration.String(),
        "max_idle_closed":      stats.MaxIdleClosed,
        "max_lifetime_closed":  stats.MaxLifetimeClosed,
    })
}
```

**Impacto:** Reduce overhead de conexiones, mejor uso de recursos

---

## 6. Prepared Statements para Queries Frecuentes

### Aplicar a queries ejecutadas 100+ veces/hora

**EJEMPLO - Login Query:**

```go
type AuthRepository struct {
    db        *sql.DB
    loginStmt *sql.Stmt  // Prepared statement
}

func NewRepository(db *sql.DB) Repository {
    repo := &pgRepository{db: db}

    // Preparar statement al inicializar
    loginStmt, err := db.Prepare(`
        SELECT account_id, email, password, phone_number, first_name, last_name,
               status, is_premium, has_finished_tutorial
        FROM account
        WHERE email = $1
    `)
    if err != nil {
        log.Fatalf("Failed to prepare login statement: %v", err)
    }
    repo.loginStmt = loginStmt

    return repo
}

func (r *pgRepository) GetUserByEmail(email string) (User, error) {
    // Usar prepared statement (más rápido)
    row := r.loginStmt.QueryRow(email)

    var user User
    var isPremium, hasFinishedTutorial dbtypes.NullBool

    err := row.Scan(
        &user.AccountID,
        &user.Email,
        // ... resto de campos
    )
    // ... resto del código
}

// Cerrar statement al finalizar
func (r *pgRepository) Close() error {
    if r.loginStmt != nil {
        return r.loginStmt.Close()
    }
    return nil
}
```

**Aplicar prepared statements a:**
- `GetUserByEmail` (login)
- `CheckAndGetIfClusterExist` (clustering)
- `HasAccountVoted` (vote checking)
- `GetUnreadCount` (notification badge)

**Impacto:** 10-20% más rápido por query

---

## 7. Cache de Queries Frecuentes (Redis)

### Para datos que cambian poco

**EJEMPLO - Categories/Subcategories:**

```go
import (
    "github.com/go-redis/redis/v8"
    "encoding/json"
    "time"
)

type CachedRepository struct {
    db    *sql.DB
    redis *redis.Client
}

func (r *CachedRepository) GetCategories() ([]Category, error) {
    ctx := context.Background()

    // 1. Intentar obtener de cache
    cacheKey := "categories:all"
    cached, err := r.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        // Cache hit
        var categories []Category
        if err := json.Unmarshal([]byte(cached), &categories); err == nil {
            return categories, nil
        }
    }

    // 2. Cache miss - obtener de DB
    query := `SELECT inca_id, name, description, icon, code
              FROM incident_categories ORDER BY name DESC`
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []Category
    for rows.Next() {
        var c Category
        if err := rows.Scan(&c.IncaId, &c.Name, &c.Description, &c.Icon, &c.Code); err != nil {
            return nil, err
        }
        categories = append(categories, c)
    }

    // 3. Guardar en cache (5 minutos)
    if jsonData, err := json.Marshal(categories); err == nil {
        r.redis.Set(ctx, cacheKey, jsonData, 5*time.Minute)
    }

    return categories, nil
}
```

**Cachear también:**
- `GetSubcategoriesByCategoryId`
- `GetInfluencerByCode` (referrals)
- Contador de clusters activos (para reels)

**Impacto:** 50-100x más rápido para datos cacheados

---

## 8. Context con Timeout en Queries

### Prevenir queries que se cuelgan

**EJEMPLO:**

```go
func (r *pgRepository) GetClustersByLocation(inputs Inputs) ([]Cluster, error) {
    // Crear context con timeout de 2 segundos
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    query := `SELECT ... FROM incident_clusters ...`

    // Usar QueryContext en lugar de Query
    rows, err := r.db.QueryContext(ctx, query, params...)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return nil, fmt.Errorf("query timeout: map loading took >2s")
        }
        return nil, err
    }
    defer rows.Close()

    // ... resto del código
}
```

**Aplicar timeout a queries críticas:**
- Map loading: 2s
- Incident details: 1s
- Login: 500ms
- Notifications: 1s

**Impacto:** Previene slow queries que bloquean conexiones

---

## 9. Logging de Slow Queries

### Para identificar problemas en producción

**EJEMPLO:**

```go
import "time"

type LoggingDB struct {
    db            *sql.DB
    slowThreshold time.Duration
}

func (ldb *LoggingDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
    start := time.Now()
    rows, err := ldb.db.Query(query, args...)
    duration := time.Since(start)

    // Log slow queries
    if duration > ldb.slowThreshold {
        log.Printf("⚠️ SLOW QUERY (%s): %s with args %v",
            duration, truncateQuery(query), args)
    }

    return rows, err
}

func truncateQuery(query string) string {
    if len(query) > 100 {
        return query[:100] + "..."
    }
    return query
}

// Usar en producción
db := &LoggingDB{
    db:            realDB,
    slowThreshold: 200 * time.Millisecond,
}
```

**Impacto:** Detectar regresiones de performance

---

## 10. Goroutine Pool para Notificaciones

### En lugar de goroutines ilimitadas

**PROBLEMA ACTUAL:**

```go
// En SaveCluster, Save, etc.
go func(accountID int64, inclID int64) {
    SaveScore(r.db, accountID, 20)
    SaveNotification(r.db, "new_cluster", accountID, inclID)
}(accountID, id)
```

**Problema:** Con 1000 reportes simultáneos → 1000 goroutines

**SOLUCIÓN - Worker Pool:**

```go
// goroutine_pool.go
type NotificationJob struct {
    Type        string
    AccountID   int64
    ReferenceID int64
    Score       uint8
}

type NotificationPool struct {
    jobs chan NotificationJob
    db   *sql.DB
}

func NewNotificationPool(db *sql.DB, workers int) *NotificationPool {
    pool := &NotificationPool{
        jobs: make(chan NotificationJob, 1000), // Buffer de 1000 jobs
        db:   db,
    }

    // Iniciar workers
    for i := 0; i < workers; i++ {
        go pool.worker(i)
    }

    return pool
}

func (p *NotificationPool) worker(id int) {
    for job := range p.jobs {
        // Procesar job
        if job.Score > 0 {
            common.SaveScore(p.db, job.AccountID, job.Score)
        }
        common.SaveNotification(p.db, job.Type, job.AccountID, job.ReferenceID)
    }
}

func (p *NotificationPool) Submit(job NotificationJob) {
    p.jobs <- job
}

// Usar en repository
var notificationPool *NotificationPool

func init() {
    notificationPool = NewNotificationPool(db, 10) // 10 workers
}

func (r *pgRepository) SaveCluster(cluster Cluster, accountID int64) (int64, error) {
    // ... INSERT ...

    // Enviar a pool en lugar de goroutine
    notificationPool.Submit(NotificationJob{
        Type:        "new_cluster",
        AccountID:   accountID,
        ReferenceID: id,
        Score:       20,
    })

    return id, nil
}
```

**Impacto:** Controla concurrencia, evita sobrecarga

---

## Resumen de Prioridades

### Implementar AHORA (Impacto Alto)
1. ✅ Eliminar subconsultas en UPDATEs (2-3x mejora)
2. ✅ Optimizar GetReel RANDOM() (10-20x mejora)
3. ✅ Separar query JSON_AGG en Profile (3-5x mejora)

### Implementar Esta Semana (Impacto Medio)
4. ⚠️ Batch inserts para notificaciones (10-50x mejora)
5. ⚠️ Connection pooling optimization (reducir overhead)
6. ⚠️ Context timeout en queries críticas (prevenir hangs)

### Implementar Próximo Mes (Impacto Bajo)
7. 📋 Prepared statements (10-20% mejora)
8. 📋 Cache Redis para datos estáticos (50-100x cache hit)
9. 📋 Logging de slow queries (monitoring)
10. 📋 Goroutine pool (control de concurrencia)

---

## Testing Recomendado

Después de cada optimización:

```bash
# 1. Benchmark antes
go test -bench=. -benchmem -benchtime=10s > before.txt

# 2. Aplicar optimización

# 3. Benchmark después
go test -bench=. -benchmem -benchtime=10s > after.txt

# 4. Comparar resultados
benchcmp before.txt after.txt
```

---

**Ver también:**
- `SQL_QUERIES_AND_INDEXES_ANALYSIS.md` - Índices requeridos
- `critical_indexes_phase1.sql` - Script SQL
