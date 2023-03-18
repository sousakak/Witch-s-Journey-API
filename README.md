# Witch's Journey Character Data
Returns data on characters from the The Journey of Elaina (魔女の旅々). I'd like to implement it like an API using repl.it in the future.

## FEATURE
This program collects data from a spreadsheet The data is encoded into a json file and displayed like api.
The program itself returns information about the characters of the Witches' Journey via API. However, the Witches' Journeys data is contained in an environment variable and cannot be seen by others.
I am planning to make the API available to everyone by using repl.it in the future.

***
## Note
This is the point I stumbled upon in writing Go: 
- When encoding from struct to json, the first letter of the variable name must be capitalized
