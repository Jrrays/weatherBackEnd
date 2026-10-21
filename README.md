# Weather Backend
This application acts as a bacekend for the [Jrrays/weatherFrontEnd](https://github.com/Jrrays/weatherFrontEnd). It fetches weather data from OpenWeatherMap, and passes it on the user frontend.

## Requirments
OpenWeatherMap API key

## Response format
{  
    "City":  
    "Temperature":  
    "Humidity":  
    "Description":  
    "Icon":  
}

## TODO
- Send wider selection of weather data to the user
- Fetch weather data based on user's location
    - User sends their location to the backend, server fetches user's local data
- Fix hardcoded allowed CORS values
- Enable HTTPS connection
- Add new features:
    - Air pollution data (OpenWeatherMap)
    - Forecasts (OpenWeatherMap)
    - Moon data (FreeAstroApi)
