# Weather Backend
This application acts as a bacekend for the [Jrrays/weatherFrontEnd](https://github.com/Jrrays/weatherFrontEnd). It fetches weather data from OpenWeatherMap, and passes it on the user frontend.

## Requirments
OpenWeatherMap API key
FreeAstroAPI API key

## Response formats
### Weather
{  
    "City":  
    "Temperature" (C):  
    "Feels Like" (C):  
    "Humidity": (%)  
    "WindSpeed" (m/s)  
    "Description":  
    "Icon":  
}
### Moon phase
{  
    "Phase name":  
    "Phase Angle" (degree)  
    "Age": (days)  
    "Distance"  (km)  
    "svg":  (Phase visual svg data)  
    "Tradotional name":  (e.g., Harvest Moon)  
    "Applies To":  (Date of the next full moon)  
}

## TODO
- Send wider selection of weather data to the user
- Fetch weather data based on user's location
    - User sends their location to the backend, server fetches user's local data
- Fix hardcoded allowed CORS values
- Enable HTTPS connection
- Add modularity/divide in packages
- Add new features:
    - Air pollution data (OpenWeatherMap)
    - Forecasts (OpenWeatherMap)
    - Moon data (FreeAstroApi)
