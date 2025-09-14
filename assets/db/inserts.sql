-- -----------------------------------------------------
-- Data for table `alertly`.`incident_categories`
-- -----------------------------------------------------
START TRANSACTION;
USE `alertly`;
INSERT INTO `alertly`.`incident_categories` (`inca_id`, `name`, `description`, `icon_uri`, `code`, `border_color`, `default_circle_range`, `max_circle_range`) VALUES (1, 'Crime', 'Incidents involving illegal activities such as theft, assault, or other offenses that threaten public safety.', 'crime', 'crime', NULL, NULL, NULL);
INSERT INTO `alertly`.`incident_categories` (`inca_id`, `name`, `description`, `icon_uri`, `code`, `border_color`, `default_circle_range`, `max_circle_range`) VALUES (2, 'Traffic Accident', 'Incidents on roadways involving vehicles, pedestrians, or both that result in collisions, injuries, or property damage.', 'traffic_accident', 'traffic_accident', NULL, NULL, NULL);
INSERT INTO `alertly`.`incident_categories` (`inca_id`, `name`, `description`, `icon_uri`, `code`, `border_color`, `default_circle_range`, `max_circle_range`) VALUES (3, 'Medical Emergency', 'Urgent situations involving sudden illness, injury, or other health crises that require immediate attention.', 'medical_emergency', 'medical_emergency', NULL, NULL, NULL);
INSERT INTO `alertly`.`incident_categories` (`inca_id`, `name`, `description`, `icon_uri`, `code`, `border_color`, `default_circle_range`, `max_circle_range`) VALUES (4, 'Fire', ' Incidents involving flames, smoke, or explosions in residential, commercial, or natural settings that pose danger to life and property.', 'fire_incident', 'fire_incident', NULL, NULL, NULL);
INSERT INTO `alertly`.`incident_categories` (`inca_id`, `name`, `description`, `icon_uri`, `code`, `border_color`, `default_circle_range`, `max_circle_range`) VALUES (5, 'Vandalism', 'Deliberate damage, defacement, or destruction of property, including graffiti or other forms of tampering with public or private assets.', 'vandalism', 'vandalism', NULL, NULL, NULL);
INSERT INTO `alertly`.`incident_categories` (`inca_id`, `name`, `description`, `icon_uri`, `code`, `border_color`, `default_circle_range`, `max_circle_range`) VALUES (6, 'Suspicious Activity', 'Observations of unusual or potentially dangerous behavior that may indicate criminal intent or cause public concern.', 'suspicious_activity', 'suspicious_activity', NULL, NULL, NULL);
INSERT INTO `alertly`.`incident_categories` (`inca_id`, `name`, `description`, `icon_uri`, `code`, `border_color`, `default_circle_range`, `max_circle_range`) VALUES (7, 'Infrastructure Issue', 'Problems or defects in public structures or services, such as damaged roads, broken streetlights, or failing utilities that affect community functionality.', 'infrastructure_issues', 'infrastructure_issues', NULL, NULL, NULL);
INSERT INTO `alertly`.`incident_categories` (`inca_id`, `name`, `description`, `icon_uri`, `code`, `border_color`, `default_circle_range`, `max_circle_range`) VALUES (8, 'Extreme Weather', 'Severe meteorological conditions, including heavy rain, snow, high winds, or other phenomena that can disrupt normal activities and cause damage.\n\n', 'extreme_weather', 'extreme_weather', NULL, NULL, NULL);
INSERT INTO `alertly`.`incident_categories` (`inca_id`, `name`, `description`, `icon_uri`, `code`, `border_color`, `default_circle_range`, `max_circle_range`) VALUES (9, 'Community Event', 'Organized gatherings, celebrations, or public meetings that bring people together for cultural, social, or civic purposes.', 'community_events', 'community_events', NULL, NULL, NULL);
INSERT INTO `alertly`.`incident_categories` (`inca_id`, `name`, `description`, `icon_uri`, `code`, `border_color`, `default_circle_range`, `max_circle_range`) VALUES (10, 'Dangerous Wildlife Sighting', 'Reports of potentially hazardous encounters with wild animals that may pose risks to humans or property.', 'dangerous_wildlife_sighting', 'dangerous_wildlife_sighting', NULL, NULL, NULL);
INSERT INTO `alertly`.`incident_categories` (`inca_id`, `name`, `description`, `icon_uri`, `code`, `border_color`, `default_circle_range`, `max_circle_range`) VALUES (11, 'Positive Actions', 'Reports of helpful or altruistic behaviors and initiatives that benefit individuals or the community, highlighting acts of kindness and civic responsibility.', 'positive_actions', 'positive_actions', NULL, NULL, NULL);
INSERT INTO `alertly`.`incident_categories` (`inca_id`, `name`, `description`, `icon_uri`, `code`, `border_color`, `default_circle_range`, `max_circle_range`) VALUES (12, 'lost_pet', 'Cases involving missing or found pets reported to keep the community informed and help increase the chances of reuniting them with their owners.', 'lost_pet', 'lost_pet', NULL, NULL, NULL);

COMMIT;


-- -----------------------------------------------------
-- Data for table `alertly`.`incident_subcategories`
-- -----------------------------------------------------
START TRANSACTION;
USE `alertly`;
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (1, 1, 'Theft', 'Incidents involving the unauthorized taking of property without the use of force.', 'crime', 'crime');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (2, 1, 'Robbery', 'Crimes where property is taken by force or intimidation from a person or place.', 'robbery', 'robbery');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (3, 1, 'Assault', 'Physical attacks or violent confrontations that result in harm or the threat of harm.', 'assault', 'assault');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (4, 1, 'Homicide', 'Cases involving the unlawful killing of an individual, including murder and attempted murder.', 'homicide', 'homicide');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (5, 1, 'Fraud', 'Incidents of deception or financial scams intended to gain property or money illicitly.', 'fraud', 'fraud');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (6, 2, 'Vehicle Collision', 'Incidents where two or more vehicles collide, often resulting in property damage or injuries.', 'single_vehicle_accident', 'single_vehicle_accident');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (7, 2, 'Pedestrian Involvement', 'Accidents that include pedestrians being struck or endangered by vehicles.', 'pedestrian_nvolvement', 'pedestrian_nvolvement');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (8, 2, 'Hit-and-Run', 'Accidents in which the driver responsible leaves the scene without providing assistance or identification.', 'hit_and_run', 'hit_and_run');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (9, 2, 'Multi-Vehicle Pileup', 'Large-scale collisions involving several vehicles, typically resulting in complex traffic disruptions.', 'traffic_accident', 'traffic_accident');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (10, 3, 'Cardiac Arrest', 'Situations where an individual''s heart stops functioning effectively, requiring immediate intervention.', 'medical_emergency', 'medical_emergency');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (11, 3, 'Stroke', 'Emergencies where blood flow to the brain is interrupted, leading to potential brain damage.', 'medical_emergency', 'medical_emergency');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (12, 3, 'Trauma/Injury', 'Accidents or incidents resulting in physical injuries or significant bodily harm.', 'medical_emergency', 'medical_emergency');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (13, 3, 'Overdose/Poisoning', 'Cases involving excessive intake of substances or exposure to toxic materials.', 'medical_emergency', 'medical_emergency');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (14, 3, 'Other Medical Emergency', 'Any urgent medical situation that does not fall under the above categories.', 'other_medical_emergency', 'other_medical_emergency');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (15, 4, 'Building Fire', 'Fires occurring in homes or residential buildings, posing immediate danger to inhabitants.', 'building_fire', 'building_fire');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (16, 4, 'Wildfire', 'Large, uncontrolled fires that spread across forests, grasslands, or rural areas.', 'wildfire', 'wildfire');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (17, 4, 'Vehicle Fire', 'Fires involving cars or other types of vehicles, often requiring specialized response.', 'vehicle_fire', 'vehicle_fire');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (18, 4, 'Other Fire Incident', 'Fire-related incidents that do not fit into the standard categories listed above.', 'fire_incident', 'fire_incident');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (19, 5, 'Graffiti', 'Unauthorized painting, tagging, or markings on public or private property.', 'graffiti', 'graffiti');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (20, 5, 'Vehicle Vandalism', 'Acts of vandalism specifically targeting vehicles, such as keying or breaking windows.', 'vehicle_vandalism', 'vehicle_vandalism');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (21, 5, 'Public Property Damage', 'Damage to public assets like parks, monuments, or street furniture.', 'public_property_damage', 'public_property_damage');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (22, 6, 'Suspicious Person', 'Observations of individuals exhibiting unusual or concerning behavior that may warrant further attention.', 'suspicious_person', 'suspicious_person');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (23, 6, 'Suspicious Vehicle', 'Vehicles acting in an odd or unexplained manner, which could indicate potential criminal activity.', 'suspicious_vehicle', 'suspicious_vehicle');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (24, 6, 'Unusual Behavior', 'Any behavior that deviates from the norm and raises concerns for safety or security.', 'unusual_behavior', 'unusual_behavior');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (25, 6, 'Other Suspicious Activity', 'Any other behavior or situation that appears out of place and may require further investigation.', 'other_suspicious_activity', 'other_suspicious_activity');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (26, 7, 'Road Damage/Potholes', 'Damage to road surfaces, such as potholes or cracks, that may affect driving safety.', 'road_damage_potholes', 'road_damage_potholes');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (27, 7, 'Streetlight/Traffic Signal Failure', 'Malfunctions or outages in public lighting or traffic control systems.', 'streetlight_traffic_signal_failure', 'streetlight_traffic_signal_failure');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (28, 7, 'Sidewalk/Pathway Damage', 'Deterioration or hazards on pedestrian pathways and sidewalks.', 'sidewalk_pathway_damage', 'sidewalk_pathway_damage');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (29, 7, 'Public Utility Issues', 'Problems affecting essential public services like water, electricity, or gas supply.', 'public_utility_issues', 'public_utility_issues');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (30, 7, 'Structural Damage', 'Damage to critical infrastructure, such as bridges, public buildings, or monuments.', 'structural_damage', 'structural_damage');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (31, 8, 'Heavy Rain/Flooding', 'Incidents of intense rainfall that result in flooding or water accumulation on streets.', 'heavy_rain_flooding', 'heavy_rain_flooding');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (32, 8, 'Snow Storm', 'Severe winter weather conditions characterized by heavy snowfall and reduced visibility.', 'snow_storm', 'snow_storm');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (33, 8, 'Icy Road Conditions', 'Hazardous conditions on roadways due to ice formation, increasing the risk of accidents.', 'icy_roads', 'icy_roads');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (34, 8, 'High Winds/Tornado', 'Events involving extremely strong winds or tornado formation that can cause significant damage.', 'high_winds_tornado', 'high_winds_tornado');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (35, 8, 'Extreme Heat', 'Periods of unusually high temperatures that may pose health risks and strain resources.', 'extreme_heat', 'extreme_heat');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (36, 8, 'Hail/Severe Storm', 'Storms that produce large hail or other severe weather phenomena leading to property damage or hazards.', 'hail_severe_storm', 'hail_severe_storm');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (37, 9, 'Festival/Fair', 'Organized community celebrations, festivals, or fairs that bring people together.', 'festival_fair', 'festival_fair');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (38, 9, 'Public Gathering/Rally', 'Large-scale meetings, rallies, or protests that involve community participation.', 'public_gathering_rally', 'public_gathering_rally');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (39, 9, 'Concert/Performance', 'Public musical or performance events, either outdoor or in community venues.', 'concert_performance', 'concert_performance');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (40, 9, 'Community Meeting/Block Party', 'Neighborhood gatherings or block parties designed to foster community interaction.', 'community_meeting_block_party', 'community_meeting_block_party');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (41, 9, 'Sporting Event', 'Organized sports matches or community athletic events that engage local residents.', 'sporting_event', 'sporting_event');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (42, 10, 'Moose', 'Sightings of moose, a common and sometimes unpredictable species in Canada.', 'moose', 'moose');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (43, 10, 'Cougar', 'Observations of cougars, which may pose a risk in certain areas.', 'cougar', 'cougar');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (44, 10, 'Bear', 'Reports of bear sightings, including grizzlies and black bears.', 'bear', 'bear');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (45, 10, 'Wolf', 'Encounters with wolves, common in some regions of Canada.', 'wolf', 'wolf');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (46, 10, 'Coyotes', 'Sightings of coyotes, which are prevalent in urban and rural settings.', 'coyotes', 'coyotes');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (47, 10, 'Bigfoot', 'Unverified sightings that could be attributed to local folklore.', 'bigfoot', 'bigfoot');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (48, 11, 'Good Samaritan Acts', 'Reports of individuals who step in to help others during emergencies or crises without expecting anything in return.', 'good_samaritan_acts', 'good_samaritan_acts');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (49, 11, 'Community Service', 'Instances of organized or spontaneous volunteer efforts that benefit the community at large.', 'community_service', 'community_service');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (50, 11, 'Environmental Stewardship', 'Actions aimed at protecting or enhancing the local environment, such as community cleanups or tree plantings.', 'environmental_stewardship', 'environmental_stewardship');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (51, 11, 'Neighborly Assistance', 'Examples of neighbors or community members helping each other with everyday tasks or in times of need.', 'neighborly_assistance', 'neighborly_assistance');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (52, 11, 'Random Acts of Kindness', 'Spontaneous and unexpected gestures that spread positivity and support throughout the community.', 'random_acts_of_kindness', 'random_acts_of_kindness');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (53, 12, 'Dog', 'Lost or found dog. Please provide breed, color, and distinctive features.', 'lost_dog', 'lost_dog');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (54, 12, 'Cat', 'Lost or found cat. Please provide details such as color, size, and any identifying marks.', 'lost_cat', 'lost_cat');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (55, 12, 'Bird', 'Lost or found bird. Include species if known and description of plumage.', 'lost_bird', 'lost_bird');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (56, 12, 'Reptile', 'Lost or found reptile. Include details about the species and any unique markings.', 'lost_reptile', 'lost_reptile');
INSERT INTO `alertly`.`incident_subcategories` 
  (`insu_id`, `inca_id`, `name`, `description`, `icon`, `code`) 
VALUES 
  (57, 12, 'Other', 'Lost or found pet that does not fit into the above categories. Please provide additional details.', 'lost_other', 'lost_other');
COMMIT;
