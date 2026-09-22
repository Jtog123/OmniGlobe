world bank api info

Yes. The page you found is mostly a high-level summary. The World Bank maintains several dedicated documentation portals, guides, and OpenAPI schemas depending on which specific dataset you want to query.

Going to rewrite the API for more convenient use


Here are the best resources for comprehensive World Bank API documentation, complete with working examples:

---

### 1. Developer Documentation Portal (Best Overall Guide)

* **URL:** [https://datahelpdesk.worldbank.org/knowledgebase/topics/125589-developer-information](https://datahelpdesk.worldbank.org/knowledgebase/topics/125589-developer-information?utm_source=gemini)
* **Why it's better:** This is the actual Data Help Desk knowledge base. It provides clear parameter breakdowns, HTTP response codes, and step-by-step examples for querying **Countries**, **Indicators**, **Topics**, and **Sources**.
* **Key pages inside:**
* **Basic Call Structures:** Explains pagination (`page`, `per_page`), sorting, and language options.
* **Indicator Queries:** Shows how to filter by income levels, regions, or custom date ranges.



---

### 2. World Bank API Basic Call Structure & Examples

* **URL:** [https://datahelpdesk.worldbank.org/knowledgebase/articles/889392-about-the-indicators-api-documentation](https://datahelpdesk.worldbank.org/knowledgebase/articles/889392-about-the-indicators-api-documentation?utm_source=gemini)
* **Why it's better:** Gives quick-reference copy-paste URL examples for the most common querying scenarios:
* Fetching data for a specific region (e.g., Sub-Saharan Africa or EU).
* Fetching data by income level (e.g., Low Income Countries).
* Requesting Most Recent Values (`mrv=1`).



---

### 3. Official OpenAPI / Swagger Specification

* **URL:** [https://api.worldbank.org/v2/swagger.json](https://www.google.com/search?q=https://api.worldbank.org/v2/swagger.json&utm_source=gemini) or [https://raw.githubusercontent.com/worldbank/worldbank-api-docs](https://www.google.com/search?q=https://github.com/worldbank/worldbank-api-docs&utm_source=gemini)
* **Why it's better:** You can load this raw Swagger JSON into tools like **Postman** or **Insomnia** (or use Go's OpenAPI generators) to inspect every endpoint, query parameter, and schema definition interactively.

---

### 4. Interactive World Bank Data Catalog & API Explorer

* **URL:** [https://datacatalog.worldbank.org](https://datacatalog.worldbank.org?utm_source=gemini)
* **Why it's better:** If you aren't sure what indicator code you need (like `NY.GDP.MKTP.CD`), searching the Data Catalog allows you to find datasets visually and reveals the exact API endpoints and indicator IDs on each dataset's detail page.

---

### Quick Cheat Sheet of Common Query Examples

To save you time digging through the docs, here are the most common API query patterns:

```http
# 1. Get metadata for ALL available indicators (500 per page)
GET https://api.worldbank.org/v2/indicator?format=json&per_page=500

# 2. Get 10 years of GDP data for all G7 countries
GET https://api.worldbank.org/v2/country/CAN;FRA;DEU;ITA;JPN;GBR;USA/indicator/NY.GDP.MKTP.CD?format=json&date=2014:2024

# 3. Get the Most Recent Value (MRV) for Population across low-income countries
GET https://api.worldbank.org/v2/incomeLevel/LIC/indicator/SP.POP.TOTL?format=json&mrv=1

# 4. Search indicators by keyword (e.g., "imports")
GET https://api.worldbank.org/v2/indicator?format=json&q=imports

```


///////////




Yes. **For your project, I think that's the right architecture once you have more than a tiny prototype.** But I'd make one distinction:

> Don't save the `Country` object simply because you called the World Bank API. Save the **normalized data your application actually owns**.

The flow I'd use is:

```text
World Bank API
      ↓
worldbank.Client
      ↓
normalize / transform
      ↓
Country / EconomicData models
      ↓
PostgreSQL
      ↓
Your Go API
      ↓
React / 3D Globe
```

The World Bank API is already designed for programmatic access to country and indicator data, including multiple indicators and date ranges. ([World Bank Data Help Desk][1])

### Why put Postgres in the middle?

Suppose your globe needs:

```text
USA
 ├── population
 ├── GDP
 ├── GDP per capita
 ├── exports
 ├── imports
 ├── military spending
 └── etc.
```

You **could** have your Go server call World Bank every time somebody clicks the USA.

But that's not what I'd build.

Instead, periodically ingest the data:

```text
                    ┌──────────────┐
                    │ World Bank  │
                    └──────┬───────┘
                           │
                       scheduled
                        ingestion
                           │
                           ▼
                    ┌──────────────┐
                    │ PostgreSQL  │
                    └──────┬───────┘
                           │
                     your queries
                           │
                           ▼
                    ┌──────────────┐
                    │   Go API    │
                    └──────┬───────┘
                           │
                           ▼
                       React UI
```

Now your application can do:

```http
GET /api/countries/USA
```

and Postgres handles the query immediately.

### But don't make `Country` contain everything

This is where I'd slightly change what you were thinking.

Don't make:

```go
type Country struct {
    Name          string
    Population    int64
    GDP           float64
    Imports       float64
    Exports       float64
    ...
}
```

because you'll eventually have **hundreds/thousands of observations**.

Instead:

```go
type Country struct {
    ID        string
    Name      string
    ISO2      string
    ISO3      string
    Latitude  float64
    Longitude float64
}
```

Then separate economic observations:

```go
type IndicatorValue struct {
    CountryID   string
    IndicatorID string
    Year        int
    Value       float64
}
```

So your database might look roughly like:

```text
countries
--------------------------------
id | iso3 | name | latitude | longitude

indicators
--------------------------------
id | code | name | unit

indicator_values
--------------------------------
country_id | indicator_id | year | value
```

Then:

```text
USA
 │
 ├── GDP ─────────── 2025 → $...
 ├── Population ──── 2025 → ...
 ├── GDP/capita ──── 2025 → ...
 └── CO2 emissions ─ 2025 → ...
```

That's much more scalable.

### And this becomes really powerful when you add Comtrade

You'll eventually have:

```text
                    YOUR DATABASE

countries ────────────────┐
                           │
indicators ───────────────┤
                           │
trade flows ──────────────┤
                           │
conflicts ────────────────┤
                           │
agreements ───────────────┤
                           │
sanctions ────────────────┘
```

Now your application isn't really a **World Bank application** anymore.

It's your own **world data API**.

The World Bank is just one upstream provider.

That's exactly the abstraction I'd want for the project you're describing.

### One important thing: don't overbuild this yet

For your first version, I'd literally do:

```text
internal/
├── country/
│   └── country.go
│
├── worldbank/
│   └── client.go
│
├── database/
│   └── postgres.go
│
└── api/
    └── handlers.go
```

Then implement one vertical slice:

**World Bank → Go → Postgres → `GET /countries/:id`**

Get that working for **population + GDP for maybe 10 countries**.

Once that works, add the rest.

That will teach you considerably more backend engineering than immediately creating a giant architecture.

[1]: https://datahelpdesk.worldbank.org/knowledgebase/articles/889392?utm_source=chatgpt.com "About the Indicators API Documentation – World Bank Data Help Desk"
