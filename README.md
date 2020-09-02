# Weather Data Server
This Software is part of [weather-data-system](https://github.com/ChristophBe/weather-data-system).
This Repository contains the code for a server written in go that can be used to manage weather data and provide it to frontend apps. 
To access and manipulate the Data this Server provides a Rest-API. To Store Data this project uses an Neo4j Graph Database. 

## Run and Build  
To create an executable just run `go build github.com/ChristophBe/weather-data-server`. After the build you can run the resulting executable. 

## Configuration
This server needs an configuration file. Its default name is `config.json` and is located in the same folder as the executable. 
With the `–config path\configfile.json` parameter you can use other config files.
To create your Configuration file you can make an copy of the `config_sample.json` provided in this repository and add the needed information of your system. 
