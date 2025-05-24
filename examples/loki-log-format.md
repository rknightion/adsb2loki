# Loki Log Format for adsb2loki

## Label and Metadata Strategy

To avoid cardinality explosions in Loki while maintaining efficient querying:

### Labels (Low Cardinality)
- **app**: Always set to "flightaware"

### Structured Metadata (Medium Cardinality)
- **hex**: Aircraft identifier (e.g., "4ca614")
- **flight**: Flight number (e.g., "EIN581")
- **category**: Aircraft category (e.g., "A5") - only included when present

Structured metadata provides indexed access without the cardinality issues of labels, making queries fast while keeping the index size manageable.

## Example Log Entry

### Labels
```
{app="flightaware"}
```

### Log Line (JSON)
```json
{
  "hex": "4ca614",
  "flight": "EIN581",
  "r": "EI-FNG",
  "t": "A333",
  "desc": "AIRBUS A-330-300",
  "alt_baro": 39950,
  "alt_geom": 40850,
  "gs": 453.7,
  "lat": 50.99928,
  "lon": -6.054611,
  "squawk": "6016",
  "emergency": "none",
  "category": "A5",
  "seen": 0.8,
  "rssi": -33.1
}
```

## Querying in Loki

With structured metadata, you can query efficiently using both metadata filters and JSON parsing:

### Using Structured Metadata (Fast)
```logql
# Find all logs for a specific aircraft
{app="flightaware", hex="4ca614"}

# Find all logs for a specific flight
{app="flightaware", flight="EIN581"}

# Find all aircraft in a specific category
{app="flightaware", category="A5"}

# Combine metadata filters
{app="flightaware", flight="EIN581", category="A5"}
```

### Using JSON Parsing (For fields not in metadata)
```logql
# Find all aircraft above certain altitude
{app="flightaware"} | json | alt_baro > 35000

# Complex queries combining metadata and JSON
{app="flightaware", category="A5"} | json | alt_baro > 35000 and gs > 400
```

This approach provides the best of both worlds:
- Fast indexed queries for common filters (hex, flight, category) via structured metadata
- Full flexibility to query any field via JSON parsing
- No cardinality explosion in the label index 