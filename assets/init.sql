--
-- page
--
define table overwrite page schemafull type any;

-- route fields
define field method on table page type string;
define field path on table page type string;
define field name on table page type string;

-- content fields
define field content on table page type object;
define field parent on table page type record; -- todo: specify allowed record types

