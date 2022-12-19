# Notes

- You can run the application with `docker-compose up`
- To change the target url for the proxy change the env variable `TARGET_URL` (I've also tested forwarded to localhost:8080)
- Regarding the proxy implementation... I tried to implement this using low level sockets but couldn't get it to work 
  so afraid I've gone down the easy route with the proxy 
- As the server is only looking for a query param i've decided to only allow GET methods
- You can send a valid request to the server with `http://localhost:8080/?name=${name}`
