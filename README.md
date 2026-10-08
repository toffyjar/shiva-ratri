# Jyotish API

A REST API for Jyotish (Vedic Astrology) calculations, built on top of the [kunjara/jyotish](https://github.com/kunjara/jyotish) library. Jyotish API  containerized and provides a wide range of Vedic astrology calculations, including planetary positions, Dashas, Yogas, and more.

Showcase : [https://kebrya.com/chart/calculator](https://kebrya.com/chart/calculator)
## Features

### Chart Calculations

- **Varga (Divisional) Charts**: D1, D2, D3, D4, D7, D9, D10, D12, D16, D20, D24, D27, D30, D40, D45, D60
- **Planetary Positions**: Detailed positions of planets, Surya (Sun), Chandra (Moon), Mangal (Mars), Buddha (Mercury), Guru (Jupiter), Shukra (Venus), Shani (Staurn), Rahu and Ketu.
- **Lagna (Ascendant) Calculations**.
- **Nakshatra Positions**: Lunar mansions.

### Advanced Calculations

#### Basic Information

- **Astangata (Combustion)**: Analysis of planets combust due to proximity to the Sun.
- **Rashi Avastha (Planetary Dignity)**: Calculation of planetary dignity in different signs.
- **Vargottama**: Identification of planets occupying the same sign in both Rashi and Navamsa charts.
- **Planetary War (Yuddha)**: Detection of planetary wars when planets are within 1 degree of each other.

#### Panchanga Elements

- **Tithi**: Lunar day calculation.
- **Nakshatra**: Determination of the current lunar mansion.
- **Yoga**: Calculation of the Yoga formed by the positions of the Sun and Moon.
- **Karana**: Half of a Tithi, used in muhurta (electional astrology).
- **Vara**: Day of the week.

#### Strength Calculations

- **Ashtakavarga**: Point-based system to analyze planetary strengths.
- **Graha Bala (Planetary Strength)**: Comprehensive strength analysis of planets.
- **Rashi Bala (Sign Strength)**: Strength of the zodiac signs.

#### Yoga Calculations

- **Dhana Yoga**: Wealth-producing combinations.
- **Mahapurusha Yoga**: Great person yogas formed by certain planetary positions.
- **Nabhasha Yoga**: Yogas based on planetary patterns.
- **Parivarthana Yoga**: Mutual exchange of signs between planets.
- **Raja Yoga**: Combinations indicating power and authority.
- **Sannyasa Yoga**: Combinations indicating renunciation.

### Additional Features

- **Automatic Timezone Detection**: Adjusts calculations based on the provided location.
- **DST Handling**: Accounts for Daylight Saving Time shifts.
- **Ayanamsha Calculations**: Adjustments for the precession of the equinoxes.
- **Dasha Calculations (Vimshottari)**: Planetary periods analysis.
- **Supported Ayanamshas**: Deluce ,Djwhalkhul ,Fagan ,Jnbhasin ,Krishnamurti ,Lahiri ,Raman ,Sassanian ,Ushashashi ,Yukteshwar

## Literature

Most of the calculations are based on the following classical texts:

- Maharishi Parashara. Brihat Parashara Hora Shastra.
- Maharishi Jaimini. Jaimini Upadesha Sutras.
- Varahamihira. Brihat Jataka.
- Varahamihira. Brihat Samhita.
- Kalyana Varma. Saravali.
- Satyacharya. Satya Jatakam.
- Kalidas. Uttara Kalamritam.
- Venkatesh Sharma. Sarvarth Chintamani.
- Mantreswara. Phaladeepika.
- Vaidyanatha Dikshita. Jataka Parijata.
- Srimad-Bhagavatam.
- Bhavishya Purana.
- Surya Siddhanta.
- Manu-Samhita.


## Installation

To install and run the Vedic Astrology API locally, follow these steps:

1. **Clone the repository**:

   ```bash
   git clone https://github.com/teal33t/jyotish-api.git
   cd jyotish-api
   ```

2. **Build and start the Docker container**:

   ```bash
   docker-compose up --build -d
   ```

3. **Compile the Swiss Ephemeris library inside the container**:

   ```bash
   docker-compose exec jyotish_api bash -c "cd swetest/src && make clean && make && chmod 777 swetest/src && chmod +x swetest/src/swetest"
   ```

4. **Verify the API is running**:

   Visit [http://localhost:9393/api/ping](http://localhost:9393/api/ping) to check if the API is up and running.

## Access control (optional)

The API is open by default. To restrict it, set these environment variables on the host (never in a committed file):

- `SR_API_KEYS`: one or more keys, comma-separated. When set, every `/api/*` request except `/api/ping` must send one of them in the `X-SR-Key` header, or it gets `401`.
- `SR_ALLOWED_IPS` (optional): IPs or CIDR ranges, comma-separated, allowed without a key (for example, your own machine for the Swagger page).

`/` and `/api/ping` stay open for health checks. The caller's address is the right-most public address in `X-Forwarded-For`, so a caller cannot fake an allowed address. Give each client its own key, so one can be revoked without touching the others.

## API Documentation

The API documentation is available via Swagger UI:

- **Swagger UI**: [http://localhost:9393/api/doc](http://localhost:9393/api/doc)

![Swagger](swagger_api_doc.png)


## API Endpoints

### Health Check

```http
GET /api/ping
```

Verifies if the API is running.

**Response**:
```json
{
    "pong": "success"
}
```

### Calculate Chart

```http
GET /api/calculate
```

Calculates chart based on provided parameters.


**Request**

```bash
http://localhost:9393/api/calculate?latitude=37.28077&longitude=49.583057&year=2023&month=12&day=25&hour=12&min=0&sec=0&time_zone=Asia/Tehran&dst_hour=0&dst_min=0&nesting=0&varga=D1,D9&infolevel=basic,panchanga,transit
```

**Query Parameters**:

| Parameter  | Type    | Required | Description                           | Example    |
|------------|---------|----------|---------------------------------------|------------|
| latitude   | float   | Yes      | Location latitude                     | 28.6139    |
| longitude  | float   | Yes      | Location longitude                    | 77.2090    |
| year       | integer | Yes      | Year for calculation                  | 2023       |
| month      | integer | Yes      | Month for calculation (1-12)          | 12         |
| day        | integer | Yes      | Day for calculation (1-31)            | 25         |
| hour       | integer | Yes      | Hour for calculation (0-23)           | 12         |
| min        | integer | Yes      | Minute for calculation (0-59)         | 0          |
| sec        | integer | Yes      | Second for calculation (0-59)         | 0          |
| time_zone  | string  | No       | Timezone for calculation             | Asia/Tehran |
| dst_hour   | integer | No       | DST hours offset                      | 0          |
| dst_min    | integer | No       | DST minutes offset                    | 0          |
| nesting    | integer | No       | Nesting level for calculations        | 0          |
| varga      | string  | No       | Varga divisions (comma-separated)     | D1,D9      |
| infolevel  | string  | No       | Info levels to include               | basic,panchanga |

**Response**:
```json

{
  "chart": {
    "user": {
      "datetime": "2023-12-25 12:00:00",
      "timezone": "+03:30",
      "longitude": 49.583057,
      "latitude": 37.28077,
      "altitude": 0
    },
    "graha": {
      "Sy": {
        "longitude": 248.983163,
        "latitude": 0.0001772,
        "speed": 1.0183709,
        "ascension": 273.4600331,
        "declination": -23.3978588,
        "rashi": 9,
        "degree": 8.98316299999999,
        "nakshatra": {
          "anga": "nakshatra",
          "key": 19,
          "ratio": 1,
          "abhijit": false,
          "left": 32.62627750000016,
          "name": "Moola",
          "pada": 3
        },
        "astangata": null,
        "rashiAvastha": "friend",
        "vargottama": false,
        "yuddha": null,
        "gocharastha": 0,
        "bhavaCharacter": "mishra",
        "tempRelation": {
          "Ch": -1,
          "Ma": 1,
          "Bu": -1,
          "Gu": -1,
          "Sk": 1,
          "Sa": 1,
          "Ra": 1,
          "Ke": 1
        },
        "relation": {
          "Ch": 0,
          "Ma": 2,
          "Bu": -1,
          "Gu": 0,
          "Sk": 0,
          "Sa": 0,
          "Ra": 0,
          "Ke": 0,
          "Sy": null
        },
        "yogakaraka": false,
        "mrityu": false,
        "pushkaraNavamsha": false,
        "pushkaraBhaga": false,
        "avastha": {
          "baladi": "kumara",
          "jagradi": "swapna",
          "deeptadi": [
            "shanta"
          ]
        },
        "dispositor": "Gu"
      },
      "Ch": {},
      "Ma": {},
      "Bu": {},
      "Gu": {},
      "Sk": {},
      "Sa": {},
      "Ra": {},
      "Ke": {},
    },
    "bhava": {
      "1": {
        "longitude": 327.9252867,
        "ascension": 330.1017467,
        "declination": -12.1928517,
        "rashi": 11,
        "degree": 27.925286700000015
      },
      "2": {},
      "3": {},
      "4": {},
      "5": {},
      "6": {},
      "7": {},
      "8": {},
      "9": {},
      "10": {},
      "11": {},
      "12": {},
    },
    "lagna": {
      "Lg": {
        "longitude": 327.9252867,
        "ascension": 330.1017467,
        "declination": -12.1928517,
        "rashi": 11,
        "degree": 27.925286700000015,
        "nakshatra": {
          "anga": "nakshatra",
          "key": 25,
          "ratio": 1,
          "abhijit": false,
          "left": 40.56034975,
          "name": "Purva Bhadrapada",
          "pada": 3
        }
      },
      "MLg": {
        "longitude": 240.7308939,
        "ascension": 238.5809717,
        "declination": -20.3009012,
        "rashi": 9,
        "degree": 0.7308939000000123
      }
    },
    "varga": {
      "D9": {
        "bhava": {
          "1": {
            "rashi": 3,
            "degree": 11.327580300000125,
            "longitude": 71.32758030000012
          },
          "2": {},
          "3": {},
          "4": {},
          "5": {},
          "6": {},
          "7": {},
          "8": {},
          "9": {},
          "10": {},
          "11": {},
          "12": {},
        },
        "graha": {
          "Sy": {
            "rashi": 3,
            "degree": 20.84846699999991,
            "speed": 1.0183709,
            "longitude": 80.84846699999991
          },
          "Ch": {},
          "Ma": {},
          "Bu": {},
          "Gu": {},
          "Sk": {},
          "Sa": {},
          "Ra": {},
          "Ke": {},
        },
        "lagna": {
          "Lg": {
            "rashi": 3,
            "degree": 11.327580300000125,
            "longitude": 71.32758030000012
          },
          "MLg": {
            "rashi": 1,
            "degree": 6.578045100000111,
            "longitude": 6.578045100000111
          }
        }
      }
    },
    "panchanga": {
      "tithi": {
        "anga": "tithi",
        "key": 14,
        "name": "Chaturdashi",
        "paksha": "shukla",
        "left": 75.43356083333303
      },
      "nakshatra": {
        "anga": "nakshatra",
        "key": 4,
        "ratio": 1,
        "abhijit": false,
        "left": 40.51648225000003,
        "name": "Rohini",
        "pada": 3
      },
      "yoga": {
        "anga": "yoga",
        "key": 23,
        "name": "Shubha",
        "left": 73.14275974999997
      },
      "vara": {
        "anga": "vara",
        "left": 81.29129302863754,
        "key": "Ch",
        "week": "1",
        "name": "Somavar"
      },
      "karana": {
        "anga": "karana",
        "key": 5,
        "name": "Gara",
        "left": 50.867121666666065
      }
    },
    "rising": {
      "Sy": [
        {
          "rising": "2023-12-24 07:30:04",
          "setting": "2023-12-24 17:40:54"
        },
        {
          "rising": "2023-12-25 07:30:31",
          "setting": "2023-12-25 17:41:27"
        },
        {
          "rising": "2023-12-26 07:30:56",
          "setting": "2023-12-26 17:42:02"
        }
      ]
    },
    "kala": {
      "hora": {
        "number": 7,
        "key": "Bu",
        "interval": 3054.6666666666665,
        "left": 47.173723264949835,
        "type": "yama",
        "end": "2023-12-25 13:26:54"
      }
    },
    "dasha": {
      "nesting": 0,
      "type": "vimshottari",
      "key": "",
      "duration": 3786834240,
      "start": "2018-01-12 22:49:09",
      "end": "2138-01-13 01:13:09",
      "periods": {
        "Ch": {
          "nesting": 1,
          "type": "mahadasha",
          "key": "Ch",
          "duration": 315569520,
          "start": "2018-01-12 22:49:09",
          "end": "2028-01-13 09:01:09"
        },
        "Ma": {},
        "Ra": {},
        "Gu": {},
        "Sa": {},
        "Bu": {},
        "Ke": {},,
        "Sk": {},,
        "Sy": {},
      }
    }
  },
  "duration_of_response": 0.12,
  "created_at": "2025-02-08 23:55:46"
}

```

### Current Time Chart

```http
GET /api/now
```

Calculates chart for the current moment.

**Query Parameters**:

| Parameter  | Type    | Required | Description        | Default    |
|------------|---------|----------|--------------------|------------|
| latitude   | float   | No       | Location latitude  | 35.7219    |
| longitude  | float   | No       | Location longitude | 51.3347    |

**Response**:
```json
{
    "chart": { /* Chart data */ },
    "duration_of_response": 0.123,
    "created_at": "2023-12-25 12:00:00"
}
```

## Error Handling

In case of errors, the API returns a JSON response with an error message and appropriate HTTP status code:

```json
{
    "status": 400,
    "message":"No route found for GET http://localhost:9393/api/not_registered_route/"
}
```


## Support

For issues, feature requests, or questions, please open an issue in the GitHub repository.

## License

This project is licensed under the GNU License V3 - see the LICENSE file for details.

The Jyotish library used in this project is licensed under the GNU General Public License version 2 or later. See the [Jyotish LICENSE](https://github.com/kunjara/jyotish/blob/master/LICENSE) for more details.

## Contact

For any inquiries, questions, or support, please contact [on Telegram](http://t.me/samanesmaeil).
