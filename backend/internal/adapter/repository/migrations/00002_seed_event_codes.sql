-- +goose Up
-- Seed the standard SAME/EAS event codes from NOAA.
-- Source: https://www.weather.gov/nwr/eventcodes
-- These are static enough to seed; use the UI "Refresh from NOAA" to update them.

INSERT INTO event_codes (code, description, category, is_warning) VALUES
-- Warnings (immediate threat to life/property)
('BZW', 'Blizzard Warning', 'Warning', true),
('CFW', 'Coastal Flood Warning', 'Warning', true),
('DSW', 'Dust Storm Warning', 'Warning', true),
('EWW', 'Extreme Wind Warning', 'Warning', true),
('FFW', 'Flash Flood Warning', 'Warning', true),
('FLW', 'Flood Warning', 'Warning', true),
('FRW', 'Fire Warning', 'Warning', true),
('FSW', 'Flash Freeze Warning', 'Warning', true),
('FZW', 'Freeze Warning', 'Warning', true),
('GLW', 'Gale Warning', 'Warning', true),
('HFW', 'Hurricane Force Wind Warning', 'Warning', true),
('HUW', 'Hurricane Warning', 'Warning', true),
('HWW', 'High Wind Warning', 'Warning', true),
('ISW', 'Ice Storm Warning', 'Warning', true),
('LAW', 'Landslide Warning', 'Warning', true),
('LEW', 'Law Enforcement Warning', 'Warning', true),
('LSW', 'Land Slide Warning', 'Warning', true),
('MAW', 'Marine Weather Warning', 'Warning', true),
('NAW', 'National Warning', 'Warning', true),
('NMW', 'Nuclear Power Plant Warning', 'Warning', true),
('RHW', 'Radiological Hazard Warning', 'Warning', true),
('SMW', 'Special Marine Warning', 'Warning', true),
('SPW', 'Shelter in Place Warning', 'Warning', true),
('SQW', 'Snow Squall Warning', 'Warning', true),
('SSW', 'Storm Surge Warning', 'Warning', true),
('SVW', 'Severe Thunderstorm Warning', 'Warning', true),
('TOE', 'State/Local 911 Telephone Outage Emergency', 'Warning', true),
('TOR', 'Tornado Warning', 'Warning', true),
('TRW', 'Tropical Storm Warning', 'Warning', true),
('TSW', 'Tsunami Warning', 'Warning', true),
('TYW', 'Typhoon Warning', 'Warning', true),
('VOW', 'Volcano Warning', 'Warning', true),
('WSW', 'Winter Storm Warning', 'Warning', true),
-- Watches
('CFA', 'Coastal Flood Watch', 'Watch', false),
('FFA', 'Flash Flood Watch', 'Watch', false),
('FLA', 'Flood Watch', 'Watch', false),
('HUA', 'Hurricane Watch', 'Watch', false),
('SVA', 'Severe Thunderstorm Watch', 'Watch', false),
('TOA', 'Tornado Watch', 'Watch', false),
('TRA', 'Tropical Storm Watch', 'Watch', false),
('TSA', 'Tsunami Watch', 'Watch', false),
-- Advisories
('CDA', 'Child Abduction Emergency', 'Advisory', false),
('MHW', 'Marine Weather Statement', 'Advisory', false),
('NOW', 'Short Term Forecast', 'Advisory', false),
('SCY', 'Hazardous Seas', 'Advisory', false),
('SDS', 'Special Weather Statement', 'Advisory', false),
('SPS', 'Special Weather Statement', 'Advisory', false),
('SRF', 'Surf Advisory', 'Advisory', false),
('STA', 'Storm Surge Watch', 'Advisory', false),
('WFA', 'Wild Fire Watch', 'Advisory', false),
-- Statements
('ADR', 'Administrative Message', 'Statement', false),
('AVA', 'Avalanche Watch', 'Statement', false),
('AVW', 'Avalanche Warning', 'Statement', true),
('BHS', 'Beach Hazards Statement', 'Statement', false),
('DBA', 'Dam Watch', 'Statement', false),
('DBW', 'Dam Break Warning', 'Statement', true),
('DEW', 'Destructive Weather', 'Statement', false),
('EVA', 'Evacuation Watch', 'Statement', false),
('EVW', 'Evacuation Warning', 'Statement', true),
('FLS', 'Flood Statement', 'Statement', false),
('FWW', 'Fire Weather Warning', 'Statement', true),
('HLS', 'Hurricane Local Statement', 'Statement', false),
('MWS', 'Marine Weather Statement', 'Statement', false),
-- Tests
('EAN', 'Emergency Action Notification', 'Test', false),
('EAT', 'Emergency Action Termination', 'Test', false),
('EVI', 'Immediate Evacuation Warning', 'Test', false),
('NPT', 'National Periodic Test', 'Test', false),
('RMT', 'Required Monthly Test', 'Test', false),
('RWT', 'Required Weekly Test', 'Test', false)
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DELETE FROM event_codes WHERE code IN (
    'BZW','CFW','DSW','EWW','FFA','FFW','FLA','FLW','FRW','FSW','FZW',
    'GLW','HFW','HUA','HUW','HWW','ISW','LAW','LEW','LSW','MAW','NAW',
    'NMW','RHW','SMW','SPW','SQW','SSW','SVA','SVW','TOA','TOE','TOR',
    'TRA','TRW','TSA','TSW','TYW','VOW','WSW',
    'CDA','CFA','MHW','NOW','SCY','SDS','SPS','SRF','STA','WFA',
    'ADR','AVA','AVW','BHS','DBA','DBW','DEW','EVA','EVW','FLS','FWW',
    'HLS','MWS',
    'EAN','EAT','EVI','NPT','RMT','RWT'
);
