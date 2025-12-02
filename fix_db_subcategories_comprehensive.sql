-- =====================================================
-- COMPREHENSIVE FIX FOR SUBCATEGORY MISMATCHES
-- Date: 2025-11-28
-- =====================================================
-- This script fixes ALL mismatches between:
-- - incident_subcategories table in DB
-- - Categories.tsx in frontend
-- - incident_clusters and incident_reports usage
-- =====================================================

-- =====================================================
-- PART 1: ADD MISSING SUBCATEGORIES IN DB
-- =====================================================

-- Add "theft" subcategory (missing in DB, exists in Categories.tsx line 19)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  1, -- crime category
  'Theft',
  'Incidents involving the unauthorized taking of property without the use of force.',
  'theft',
  'crime.png',
  24
);

-- Add "vehicle_collision" subcategory (missing in DB, exists in Categories.tsx line 74)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  2, -- traffic_accident category
  'Vehicle Collision',
  'Incidents where two or more vehicles collide, often resulting in property damage or injuries.',
  'vehicle_collision',
  'single_vehicle_accident.png',
  24
);

-- Add "residential_fire" subcategory (missing in DB, exists in Categories.tsx line 174)
-- This is CRITICAL - TFS scraper uses this code!
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  4, -- fire_incident category
  'Residential Fire',
  'Fires occurring in homes or residential buildings, posing immediate danger to inhabitants.',
  'residential_fire',
  'building_fire.png',
  48
);

-- Add other missing subcategories from Categories.tsx
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  3, -- medical_emergency category
  'Cardiac Arrest',
  'Situations where an individual\'s heart stops functioning effectively, requiring immediate intervention.',
  'cardiac_arrest',
  'medical_emergency.png',
  24
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  3, -- medical_emergency category
  'Stroke',
  'Emergencies where blood flow to the brain is interrupted, leading to potential brain damage.',
  'stroke',
  'medical_emergency.png',
  24
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  3, -- medical_emergency category
  'Trauma/Injury',
  'Accidents or incidents resulting in physical injuries or significant bodily harm.',
  'trauma_Injury',
  'medical_emergency.png',
  24
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  3, -- medical_emergency category
  'Overdose/Poisoning',
  'Cases involving excessive intake of substances or exposure to toxic materials.',
  'overdose_poisoning',
  'medical_emergency.png',
  24
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  2, -- traffic_accident category
  'Multi-Vehicle Pileup',
  'Large-scale collisions involving several vehicles, typically resulting in complex traffic disruptions.',
  'multi_vehicle_ileup',
  'traffic_accident.png',
  24
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  11, -- positive_actions category
  'Random Acts of Kindness',
  'Spontaneous and unexpected gestures that spread positivity and support throughout the community.',
  'random_acts_of_kindness',
  'random_acts_of_kindness.png',
  168
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  11, -- positive_actions category
  'Good Samaritan Acts',
  'Reports of individuals who step in to help others during emergencies or crises without expecting anything in return.',
  'good_samaritan_acts',
  'good_samaritan_acts.png',
  168
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  9, -- community_events category
  'Festival/Fair',
  'Organized community celebrations, festivals, or fairs that bring people together.',
  'festival_fair',
  'festival_fair.png',
  72
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  9, -- community_events category
  'Public Gathering/Rally',
  'Large-scale meetings, rallies, or protests that involve community participation.',
  'public_gathering_rally',
  'public_gathering_rally.png',
  72
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  10, -- dangerous_wildlife_sighting category
  'Moose',
  'Sightings of moose, a common and sometimes unpredictable species in Canada.',
  'moose',
  'moose.png',
  48
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  10, -- dangerous_wildlife_sighting category
  'Bear',
  'Reports of bear sightings, including grizzlies and black bears.',
  'bear',
  'bear.png',
  48
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  10, -- dangerous_wildlife_sighting category
  'Coyotes',
  'Sightings of coyotes, which are prevalent in urban and rural settings.',
  'coyotes',
  'coyotes.png',
  48
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  12, -- lost_pet category
  'Lost Dog',
  'Lost or found dog. Please provide breed, color, and distinctive features.',
  'lost_dog',
  'lost_dog.png',
  168
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  12, -- lost_pet category
  'Lost Reptile',
  'Lost or found reptile. Include details about the species and any unique markings.',
  'lost_reptile',
  'lost_reptile.png',
  168
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  8, -- extreme_weather category
  'Icy Road Conditions',
  'Hazardous conditions on roadways due to ice formation, increasing the risk of accidents.',
  'icy_roads',
  'icy_roads.png',
  48
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  8, -- extreme_weather category
  'Snow Storm',
  'Severe winter weather conditions characterized by heavy snowfall and reduced visibility.',
  'snow_storm',
  'snow_storm.png',
  48
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  8, -- extreme_weather category
  'Heavy Rain/Flooding',
  'Incidents of intense rainfall that result in flooding or water accumulation on streets.',
  'heavy_rain_flooding',
  'heavy_rain_flooding.png',
  48
);

INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  7, -- infrastructure_issues category
  'Streetlight/Traffic Signal Failure',
  'Malfunctions or outages in public lighting or traffic control systems.',
  'streetlight_traffic_signal_failure',
  'streetlight_traffic_signal_failure.png',
  72
);

-- =====================================================
-- PART 2: DELETE INVALID DUPLICATES FROM DB
-- =====================================================

-- Delete invalid "crime" subcategory (category being used as subcategory)
DELETE FROM incident_subcategories
WHERE code = 'crime' AND inca_id = 1;

-- Delete invalid category-named subcategories
DELETE FROM incident_subcategories
WHERE code IN ('fire_incident', 'traffic_accident', 'medical_emergency')
  AND code = (SELECT code FROM incident_categories WHERE inca_id = incident_subcategories.inca_id);

-- =====================================================
-- PART 3: UPDATE incident_clusters TO USE VALID CODES
-- =====================================================

-- Update any remaining invalid subcategory codes in incident_clusters
UPDATE incident_clusters
SET subcategory_code = 'other_fire_incident',
    subcategory_name = 'Other Fire Incident'
WHERE subcategory_code = 'fire_incident';

UPDATE incident_clusters
SET subcategory_code = 'single_vehicle_accident',
    subcategory_name = 'Single Vehicle Accident'
WHERE subcategory_code = 'traffic_accident';

UPDATE incident_clusters
SET subcategory_code = 'other_medical_emergency',
    subcategory_name = 'Other Medical Emergency'
WHERE subcategory_code = 'medical_emergency';

UPDATE incident_clusters
SET subcategory_code = 'assault',
    subcategory_name = 'Assault'
WHERE subcategory_code = 'crime';

-- =====================================================
-- PART 4: UPDATE incident_reports TO USE VALID CODES
-- =====================================================

UPDATE incident_reports
SET subcategory_code = 'other_fire_incident'
WHERE subcategory_code = 'fire_incident';

UPDATE incident_reports
SET subcategory_code = 'single_vehicle_accident'
WHERE subcategory_code = 'traffic_accident';

UPDATE incident_reports
SET subcategory_code = 'other_medical_emergency'
WHERE subcategory_code = 'medical_emergency';

UPDATE incident_reports
SET subcategory_code = 'assault'
WHERE subcategory_code = 'crime';

-- =====================================================
-- PART 5: VERIFICATION QUERIES
-- =====================================================

-- Verify all subcategory codes in incident_clusters exist in incident_subcategories
SELECT 'Invalid subcategories in incident_clusters:' as status;
SELECT DISTINCT c.subcategory_code, c.category_code, COUNT(*) as count
FROM incident_clusters c
LEFT JOIN incident_subcategories s ON c.subcategory_code = s.code
WHERE s.code IS NULL
GROUP BY c.subcategory_code, c.category_code;

-- Verify all subcategory codes in incident_reports exist in incident_subcategories
SELECT 'Invalid subcategories in incident_reports:' as status;
SELECT DISTINCT r.subcategory_code, COUNT(*) as count
FROM incident_reports r
LEFT JOIN incident_subcategories s ON r.subcategory_code = s.code
WHERE s.code IS NULL
GROUP BY r.subcategory_code;

-- Show all valid subcategory codes with their categories
SELECT 'Valid subcategory codes in database:' as status;
SELECT s.code as subcategory_code, c.code as category_code, s.name
FROM incident_subcategories s
JOIN incident_categories c ON s.inca_id = c.inca_id
ORDER BY c.code, s.code;

SELECT 'Fix completed!' as status;
