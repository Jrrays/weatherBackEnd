# Weather Backend
This application acts as a bacekend for the Weather App. It fetches weather data from OpenWeatherMap, and passes it on the user frontend.

## Requirments
OpenWeatherMap API key

## Response format
{  
    "city": Oulu  
    "temperature_celsius":  
    "humidity_percent":  
    "description":  
}

## TODO
- Send wider selection of weather data to the user
- Fetch weather data based on user's location
    - User sends their location to the backend, server fetches user's local data
- Fix hardcoded allowed CORS values
- Add new features:
    - Air pollution data (OpenWeatherMap)
    - Forecasts (OpenWeatherMap)
    - Moon data (FreeAstroApi)