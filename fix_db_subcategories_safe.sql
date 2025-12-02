-- =====================================================
-- SAFE FIX FOR SUBCATEGORY MISMATCHES
-- Date: 2025-11-28
-- =====================================================
-- This script ONLY ADDS missing subcategories
-- Does NOT delete anything to avoid foreign key issues
-- =====================================================

-- =====================================================
-- ADD MISSING SUBCATEGORIES FROM Categories.tsx
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
-- *** CRITICAL - TFS scraper uses this code! ***
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  4, -- fire_incident category
  'Residential Fire',
  'Fires occurring in homes or residential buildings, posing immediate danger to inhabitants.',
  'residential_fire',
  'building_fire.png',
  48
);

-- Add "cardiac_arrest" (Categories.tsx line 120)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  3, -- medical_emergency category
  'Cardiac Arrest',
  'Situations where an individual\'s heart stops functioning effectively, requiring immediate intervention.',
  'cardiac_arrest',
  'medical_emergency.png',
  24
);

-- Add "stroke" (Categories.tsx line 128)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  3, -- medical_emergency category
  'Stroke',
  'Emergencies where blood flow to the brain is interrupted, leading to potential brain damage.',
  'stroke',
  'medical_emergency.png',
  24
);

-- Add "trauma_Injury" (Categories.tsx line 136)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  3, -- medical_emergency category
  'Trauma/Injury',
  'Accidents or incidents resulting in physical injuries or significant bodily harm.',
  'trauma_Injury',
  'medical_emergency.png',
  24
);

-- Add "overdose_poisoning" (Categories.tsx line 144)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  3, -- medical_emergency category
  'Overdose/Poisoning',
  'Cases involving excessive intake of substances or exposure to toxic materials.',
  'overdose_poisoning',
  'medical_emergency.png',
  24
);

-- Add "multi_vehicle_ileup" (Categories.tsx line 98)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  2, -- traffic_accident category
  'Multi-Vehicle Pileup',
  'Large-scale collisions involving several vehicles, typically resulting in complex traffic disruptions.',
  'multi_vehicle_ileup',
  'traffic_accident.png',
  24
);

-- Add "random_acts_of_kindness" (Categories.tsx line 570)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  11, -- positive_actions category
  'Random Acts of Kindness',
  'Spontaneous and unexpected gestures that spread positivity and support throughout the community.',
  'random_acts_of_kindness',
  'random_acts_of_kindness.png',
  168
);

-- Add "good_samaritan_acts" (Categories.tsx line 538)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  11, -- positive_actions category
  'Good Samaritan Acts',
  'Reports of individuals who step in to help others during emergencies or crises without expecting anything in return.',
  'good_samaritan_acts',
  'good_samaritan_acts.png',
  168
);

-- Add "festival_fair" (Categories.tsx line 421)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  9, -- community_events category
  'Festival/Fair',
  'Organized community celebrations, festivals, or fairs that bring people together.',
  'festival_fair',
  'festival_fair.png',
  72
);

-- Add "public_gathering_rally" (Categories.tsx line 429)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  9, -- community_events category
  'Public Gathering/Rally',
  'Large-scale meetings, rallies, or protests that involve community participation.',
  'public_gathering_rally',
  'public_gathering_rally.png',
  72
);

-- Add "moose" (Categories.tsx line 475)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  10, -- dangerous_wildlife_sighting category
  'Moose',
  'Sightings of moose, a common and sometimes unpredictable species in Canada.',
  'moose',
  'moose.png',
  48
);

-- Add "bear" (Categories.tsx line 491)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  10, -- dangerous_wildlife_sighting category
  'Bear',
  'Reports of bear sightings, including grizzlies and black bears.',
  'bear',
  'bear.png',
  48
);

-- Add "coyotes" (Categories.tsx line 507)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  10, -- dangerous_wildlife_sighting category
  'Coyotes',
  'Sightings of coyotes, which are prevalent in urban and rural settings.',
  'coyotes',
  'coyotes.png',
  48
);

-- Add "lost_dog" (Categories.tsx line 593)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  12, -- lost_pet category
  'Lost Dog',
  'Lost or found dog. Please provide breed, color, and distinctive features.',
  'lost_dog',
  'lost_dog.png',
  168
);

-- Add "lost_reptile" (Categories.tsx line 617)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  12, -- lost_pet category
  'Lost Reptile',
  'Lost or found reptile. Include details about the species and any unique markings.',
  'lost_reptile',
  'lost_reptile.png',
  168
);

-- Add "icy_roads" (Categories.tsx line 374)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  8, -- extreme_weather category
  'Icy Road Conditions',
  'Hazardous conditions on roadways due to ice formation, increasing the risk of accidents.',
  'icy_roads',
  'icy_roads.png',
  48
);

-- Add "snow_storm" (Categories.tsx line 366)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  8, -- extreme_weather category
  'Snow Storm',
  'Severe winter weather conditions characterized by heavy snowfall and reduced visibility.',
  'snow_storm',
  'snow_storm.png',
  48
);

-- Add "heavy_rain_flooding" (Categories.tsx line 358)
INSERT IGNORE INTO incident_subcategories (inca_id, name, description, code, icon, default_duration_hours)
VALUES (
  8, -- extreme_weather category
  'Heavy Rain/Flooding',
  'Incidents of intense rainfall that result in flooding or water accumulation on streets.',
  'heavy_rain_flooding',
  'heavy_rain_flooding.png',
  48
);

-- Add "streetlight_traffic_signal_failure" (Categories.tsx line 312)
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
-- VERIFICATION: Show all subcategory codes
-- =====================================================
SELECT 'All subcategory codes in database after fix:' as status;
SELECT s.code as subcategory_code, c.code as category_code, s.name
FROM incident_subcategories s
JOIN incident_categories c ON s.inca_id = c.inca_id
ORDER BY c.code, s.code;

SELECT 'Fix completed successfully!' as status;
