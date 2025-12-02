-- Add "other_fire_incident" subcategory (Categories.tsx line 198)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  4, -- fire_incident category
  'Other Fire Incident',
  'Fire-related incidents that do not fit into the standard categories listed above.',
  'other_fire_incident',
  'fire_incident.png',
  48
);

SELECT 'Subcategory other_fire_incident added successfully!' as status;
