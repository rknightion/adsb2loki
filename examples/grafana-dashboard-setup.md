# Grafana ADS-B Dashboard Setup Guide

This guide will help you set up the comprehensive ADS-B flight tracking dashboard for Grafana.

## Prerequisites

- Grafana 9.0 or later
- Loki data source configured in Grafana
- adsb2loki running and pushing data to Loki

## Dashboard Features

The dashboard includes 15 different panels showcasing various aspects of your ADS-B data:

### Key Metrics (Top Row)
1. **Active Aircraft** - Total number of aircraft currently being tracked
2. **Emergency Aircraft** - Aircraft with active emergency codes
3. **Average Altitude** - Average altitude of all tracked aircraft
4. **Average Ground Speed** - Average speed of all aircraft
5. **Maximum Altitude** - Highest altitude observed
6. **Average Signal Strength** - Average RSSI across all aircraft

### Time Series Visualizations
- **Aircraft Count Over Time** - Shows total, ADS-B, and MLAT aircraft trends
- **Altitude Distribution Over Time** - Displays average, maximum, and minimum altitudes
- **Wind Speed at Aircraft Locations** - Tracks wind conditions
- **Aircraft Temperature Readings** - Shows OAT and TAT measurements

### Data Analysis Panels
- **Speed vs Altitude Heatmap** - Visualizes the relationship between speed and altitude
- **Top 20 Aircraft by Signal Strength** - Table showing best received aircraft
- **Aircraft Types Distribution** - Pie chart of aircraft types (B738, A320, etc.)
- **Average Speed by Aircraft Category** - Bar chart comparing speeds by category
- **Emergency & Special Squawk Codes** - Table monitoring emergency situations

## Installation Steps

### Method 1: Import via Grafana UI

1. Open Grafana and navigate to **Dashboards** → **Import**
2. Copy the contents of `examples/grafana-adsb-dashboard.json`
3. Paste into the "Import via panel json" text area
4. Click **Load**
5. Select your Loki data source from the dropdown
6. Click **Import**

### Method 2: Import via API

```bash
# Set your Grafana URL and API key
GRAFANA_URL="http://localhost:3000"
GRAFANA_API_KEY="your-api-key-here"

# Import the dashboard
curl -X POST \
  -H "Authorization: Bearer $GRAFANA_API_KEY" \
  -H "Content-Type: application/json" \
  -d @examples/grafana-adsb-dashboard.json \
  "$GRAFANA_URL/api/dashboards/db"
```

### Method 3: Using Grafana Provisioning

1. Copy the dashboard JSON to your Grafana provisioning directory:
```bash
cp examples/grafana-adsb-dashboard.json /etc/grafana/provisioning/dashboards/
```

2. Create a provisioning configuration file `/etc/grafana/provisioning/dashboards/adsb.yaml`:
```yaml
apiVersion: 1

providers:
  - name: 'ADS-B Dashboards'
    orgId: 1
    folder: 'Aviation'
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /etc/grafana/provisioning/dashboards
```

3. Restart Grafana

## Configuration

### Data Source Selection
The dashboard uses a template variable `$datasource` to select the Loki data source. This allows you to easily switch between different Loki instances if needed.

### Time Range
- Default: Last 1 hour
- Recommended refresh interval: 10 seconds
- You can adjust the time range using Grafana's time picker

### Panel-Specific Settings

#### Heatmap Panel (Speed vs Altitude)
- X-axis: Ground Speed (knots)
- Y-axis: Altitude (feet)
- Color scale: Turbo (exponential)
- Adjust bucket counts if needed for your data density

#### Table Panels
- Click column headers to sort
- Use the search box to filter results
- Export data using the panel menu

## Customization

### Adding Filters
To filter data for specific aircraft or conditions, modify the LogQL queries. For example:

```logql
# Filter for specific aircraft type
{app="flightaware"} | json | t="B738"

# Filter for aircraft above certain altitude
{app="flightaware"} | json | alt_baro > 30000

# Filter for specific registration
{app="flightaware", hex="4ca614"}
```

### Creating Alerts

You can create alerts based on the dashboard panels:

1. Edit a panel (e.g., Emergency Aircraft)
2. Go to the Alert tab
3. Create alert conditions:
   - Alert when emergency aircraft count > 0
   - Alert when average altitude drops below threshold
   - Alert when signal strength degrades

Example alert rule:
```yaml
- alert: EmergencySquawk
  expr: count(count by (hex) (rate({app="flightaware"} | json | emergency != "none" | __error__="" [5m]))) > 0
  for: 1m
  annotations:
    summary: "Aircraft broadcasting emergency squawk code"
```

## Performance Tips

1. **Adjust Time Range**: Shorter time ranges will load faster
2. **Use Structured Metadata**: The dashboard leverages structured metadata for efficient filtering
3. **Index Important Fields**: Consider using Loki's index labels for frequently queried fields
4. **Optimize Queries**: Use `__error__=""` to filter out parsing errors
5. **Limit Table Results**: The table panels use `topk()` to limit results

## Troubleshooting

### No Data Showing
1. Verify adsb2loki is running: `docker logs adsb2loki`
2. Check Loki is receiving data: `{app="flightaware"} | json`
3. Ensure time range includes recent data
4. Verify structured metadata is being sent

### Slow Performance
1. Reduce time range
2. Increase dashboard refresh interval
3. Check Loki resource usage
4. Consider adding more specific label selectors

### Parsing Errors
If you see parsing errors in queries:
1. Check that adsb2loki is sending valid JSON
2. Verify field names match your data
3. Use `| __error__=""` to filter out problematic entries

## Advanced Usage

### Geospatial Visualization
While this dashboard doesn't include a map panel, you can add one:

1. Install the Grafana Geomap panel plugin
2. Create a new panel with this query:
```logql
{app="flightaware"} | json | lat > 0 | line_format `{"lat":{{.lat}},"lon":{{.lon}},"hex":"{{.hex}}","flight":"{{.flight}}","alt":{{.alt_baro}}}`
```

### Custom Metrics
Create custom metrics using LogQL:

```logql
# Aircraft density by altitude band
sum by (alt_band) (
  count_over_time({app="flightaware"} 
  | json 
  | alt_baro > 0 
  | label_format alt_band="{{if le .alt_baro 10000}}0-10k{{else if le .alt_baro 20000}}10-20k{{else if le .alt_baro 30000}}20-30k{{else}}30k+{{end}}"
  [$__interval])
)
```

## Support

For issues or questions:
1. Check the adsb2loki logs
2. Verify your Loki queries in Explore view
3. Consult the Grafana and Loki documentation
4. Open an issue on the adsb2loki GitHub repository 