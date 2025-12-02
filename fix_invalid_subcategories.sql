-- Fix Invalid Subcategory Codes in incident_clusters
-- These invalid codes cause "icon_uri undefined" error in the frontend

-- Fix "fire_incident" → "other_fire_incident" (5 incidents)
UPDATE incident_clusters
SET subcategory_code = 'other_fire_incident',
    subcategory_name = 'other_fire_incident'
WHERE account_id = 1
  AND subcategory_code = 'fire_incident';

-- Fix "traffic_accident" → "single_vehicle_accident" (1 incident)
UPDATE incident_clusters
SET subcategory_code = 'single_vehicle_accident',
    subcategory_name = 'single_vehicle_accident'
WHERE account_id = 1
  AND subcategory_code = 'traffic_accident';

-- Fix "medical_emergency" → "other_medical_emergency" (2 incidents)
UPDATE incident_clusters
SET subcategory_code = 'other_medical_emergency',
    subcategory_name = 'other_medical_emergency'
WHERE account_id = 1
  AND subcategory_code = 'medical_emergency';

-- Fix "vehicle_collision" → "single_vehicle_accident" (12 incidents)
UPDATE incident_clusters
SET subcategory_code = 'single_vehicle_accident',
    subcategory_name = 'single_vehicle_accident'
WHERE account_id = 1
  AND subcategory_code = 'vehicle_collision';

-- Also fix incident_reports table
UPDATE incident_reports
SET subcategory_code = 'other_fire_incident'
WHERE account_id = 1
  AND subcategory_code = 'fire_incident';

UPDATE incident_reports
SET subcategory_code = 'single_vehicle_accident'
WHERE account_id = 1
  AND subcategory_code = 'traffic_accident';

UPDATE incident_reports
SET subcategory_code = 'other_medical_emergency'
WHERE account_id = 1
  AND subcategory_code = 'medical_emergency';

UPDATE incident_reports
SET subcategory_code = 'single_vehicle_accident'
WHERE account_id = 1
  AND subcategory_code = 'vehicle_collision';

-- Verify the fix
SELECT 'Invalid subcategories remaining:' as status;
SELECT subcategory_code, COUNT(*) as total
FROM incident_clusters
WHERE account_id = 1
  AND subcategory_code IN ('fire_incident', 'traffic_accident', 'medical_emergency', 'vehicle_collision')
GROUP BY subcategory_code;
