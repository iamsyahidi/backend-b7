# Go Zoom Starter to Manage Meetings

## Getting Started

1. Create a user-level Oauth application in the [Zoom Marketplace](https://marketplace.zoom.us/). Please refer to [public documentation](https://developers.zoom.us/docs/integrations/create/) for further instructions on doing so.
    * Keep the **Client ID** and **Client Secret** handy as we will need those soon.
2. Add the following scopes:
    * _meeting:read:meeting_
    * _meeting:write:meeting_
    * _meeting:update:meeting_
    * _meeting:delete:meeting_
    * _user:read:user_
    * _user:update:user_
    * _user:delete:user_

    **Note**: If you wish to add additional API endpoints to this application, you may need to add additional scopes.

Follow these steps to set up and run the application:

### Local Setup

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/iamsyahidi/backend-b7.git
   cd backend-b7
   run `go mod tidy -v ; go mod download` to download dependencies.
   ```
2. **Database Setup**:
   - Create a PostgreSQL database.
   - Set environment variables for your database and zoom API credentials as shown in the .env.example file.
3. **Ngrok Setup**:
    run `ngrok http 3002` to create a tunnel to port 3002 and use the _forwarding URL_ for this value. Also add this value in your marketplace app configuration for the **App Credentials** -> **Redirect URL for OAuth** and **Add allow lists** text input values.
4. **Run the Application**:
   - Directly with Go:
   ```bash
   go run main.go
   ```
   - Using Docker Compose:
   ```bash
   docker-compose up -d --build --force-recreate
   ```

   Head to your marketplace Oauth application and click on the **Activation** tab. Either click the **Add** button or copy your **Add URL** into a new tab.

   If everything above was setup correctly, you should see a Zoom Oauth Allow page! Click **Allow** to authorize the Zoom Oauth handshake. 



This server provides the following APIs:

| GET       | /v1/redirect             | Zoom Oauth Base                                      |
|-----------|--------------------------|------------------------------------------------------|
| GET       | /v1/health               | Healt check                                          |
| GET       | /v1/meets                | List meetings                                        |
| POST      | /v1/meets                | Create meeting                                       |
| PUT       | /v1/meets/:meeting_id    | Update meeting                                       |
| DELETE    | /v1/meets/:meeting_id    | Delete meeting                                       |

