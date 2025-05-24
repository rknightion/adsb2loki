# ADS-B Dashboard Preview

This document provides a visual description of the ADS-B Flight Tracking Dashboard panels.

## Dashboard Layout

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                          ADS-B Flight Tracking Dashboard                                 │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┬─────────────────┤
│   Active    │  Emergency  │   Average   │   Average   │   Maximum   │  Average Signal │
│  Aircraft   │  Aircraft   │  Altitude   │Ground Speed │  Altitude   │    Strength     │
│     20      │      0      │  28,543 ft  │  385 knots  │  43,000 ft  │    -28.5 dB     │
├─────────────┴─────────────┴─────────────┴─────────────┴─────────────┴─────────────────┤
│                                                                                         │
│                    Aircraft Count Over Time                 Altitude Distribution      │
│  ┌─────────────────────────────────────────┐  ┌─────────────────────────────────────┐ │
│  │     📈 Total: 20                         │  │     📊 Max: 43,000 ft               │ │
│  │     📈 ADS-B: 18                         │  │     📊 Avg: 28,543 ft               │ │
│  │     📈 MLAT: 2                           │  │     📊 Min: 1,350 ft                │ │
│  └─────────────────────────────────────────┘  └─────────────────────────────────────┘ │
│                                                                                         │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                         │
│  Speed vs Altitude Heatmap          Top 20 Aircraft by Signal Strength                 │
│  ┌─────────────────────┐  ┌───────────────────────────────────────────────────────┐  │
│  │ 🟦🟦🟦🟦🟦🟦🟦🟦 │  │ Hex    Flight  Reg     Type  Alt    Speed  RSSI   │  │
│  │ 🟦🟩🟩🟩🟩🟩🟦🟦 │  │ 4080cb NJU426A G-NJAF  F2TH  38000  458.7  -19.7  │  │
│  │ 🟦🟩🟨🟨🟨🟩🟦🟦 │  │ 406a05 BAW1B   G-XLEF  A388  31900  415.0  -26.9  │  │
│  │ 🟦🟩🟨🟧🟨🟩🟦🟦 │  │ 4cae6c RYR2853 EI-IJY  B38M  37000  442.2  -21.3  │  │
│  └─────────────────────┘  └───────────────────────────────────────────────────────┘  │
│                                                                                         │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                         │
│  Aircraft Types Distribution    Average Speed by Category    Wind Speed Over Time      │
│  ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────────────┐  │
│  │    B738 (35%)  🥧  │  │ A1: ▓▓▓ 156 kts     │  │ 📈 Avg Wind: 45 knots       │  │
│  │    A320 (25%)  🥧  │  │ A2: ▓▓▓▓▓ 285 kts   │  │ 📈 Max Wind: 67 knots       │  │
│  │    B77W (15%)  🥧  │  │ A3: ▓▓▓▓▓▓▓ 385 kts │  │                             │  │
│  │    Other (25%) 🥧  │  │ A5: ▓▓▓▓▓▓▓▓ 445 kt │  │                             │  │
│  └─────────────────────┘  └─────────────────────┘  └─────────────────────────────┘  │
│                                                                                         │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                         │
│  Emergency & Special Squawk Codes              Aircraft Temperature Readings            │
│  ┌─────────────────────────────────────┐  ┌─────────────────────────────────────────┐ │
│  │ Time  Hex   Flight Squawk Emergency │  │ 📈 Outside Air Temp (OAT): -45°C        │ │
│  │ 14:23 4cae6c RYR2853 7501  none     │  │ 📈 Total Air Temp (TAT): -20°C          │ │
│  │ 14:15 406669 EZY21GE 7700  hijack   │  │                                         │ │
│  └─────────────────────────────────────┘  └─────────────────────────────────────────┘ │
│                                                                                         │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

## Panel Descriptions

### Row 1: Key Metrics (Stats Panels)
- **Active Aircraft**: Shows real-time count with color thresholds (green < 50, yellow < 100, orange < 200, red ≥ 200)
- **Emergency Aircraft**: Highlights any aircraft with emergency squawk codes (red background if > 0)
- **Average Altitude**: Color-coded by altitude bands (blue < 10k, green < 25k, yellow < 35k, orange ≥ 35k)
- **Average Ground Speed**: Shows knots with thresholds (green < 300, yellow < 400, orange < 500, red ≥ 500)
- **Maximum Altitude**: Purple value display showing the highest aircraft
- **Average Signal Strength**: RSSI in dB (red < -40, orange < -30, yellow < -20, green ≥ -20)

### Row 2: Time Series Analysis
- **Aircraft Count Over Time**: Multi-line graph showing total, ADS-B, and MLAT aircraft trends
- **Altitude Distribution**: Shows min/avg/max altitude bands with gradient coloring

### Row 3: Data Visualization
- **Speed vs Altitude Heatmap**: 2D heatmap with speed on X-axis, altitude on Y-axis, using Turbo color scheme
- **Top 20 Aircraft Table**: Sortable table with RSSI color coding and unit formatting

### Row 4: Statistical Analysis
- **Aircraft Types Pie Chart**: Donut chart showing distribution of aircraft models
- **Average Speed Bar Chart**: Horizontal bars comparing speeds across aircraft categories (A1-A5)
- **Wind Speed Time Series**: Dual-line graph showing average and maximum wind speeds

### Row 5: Monitoring & Environmental
- **Emergency Squawk Table**: Real-time monitoring of 7500/7600/7700 codes with timestamp
- **Temperature Scatter Plot**: Points showing OAT (blue) and TAT (red) readings over time

## Interactive Features

1. **Hover Details**: All panels show detailed tooltips on hover
2. **Click to Filter**: Click on legend items to show/hide series
3. **Zoom**: Time series panels support click-and-drag zoom
4. **Sort**: Table columns are sortable
5. **Export**: Each panel can export data via the panel menu
6. **Full Screen**: Double-click any panel for full-screen view

## Color Schemes

- **Thresholds**: Green → Yellow → Orange → Red for increasing values
- **Gradients**: Continuous color scales for heatmaps and altitude
- **Categories**: Classic palette for distinct categories
- **Temperature**: Blue for cold (OAT), Red for warm (TAT)

## Refresh & Updates

- Auto-refresh: Every 10 seconds
- Live data: Trailing 1-hour window by default
- Instant queries: Stats update immediately
- Time series: Smooth interpolation between points 