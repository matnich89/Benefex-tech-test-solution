# My Thinking Process on the Project Specs

## Quick Overview

When I first looked at the specs, it was pretty clear that we needed a setup with distinct parts. Since there's a message broker involved, I figured we're talking microservices. So, I set up a monorepo with three microservices:

- **API**: This one includes the starter code plus the changes I made.
- **Distribution Service**: Takes care of sending out CDs and vinyls.
- **Communication Service**: Handles letting fans know about new releases.

## Initial Planning

Before jumping into coding, I spent a bit of time mapping out my approach. Here's what I had in mind:

1. Set up the monorepo with the three services, each with its own Dockerfile and a docker-compose file.
2. Work on the RabbitMQ stuff that more than one service would need.
3. Update the API, so it can send messages to RabbitMQ.
4. Get basic message handling up and running in both the communication and distribution services.
5. Dive into the specific business logic for both services.
6. Do a final check and clean up. Unfortunately, time was a bit tight, so there might be a few loose ends.

## Assumptions and Implementation Details

### Fans Database, Email Client, and Distribution Clients

In the development of this project, I made a few assumptions that shaped the implementation strategy, particularly concerning external dependencies like the fans database, email client, and distribution clients. Here's a brief overview:


- **Fans Database**: I assumed that the database holding information about the fans (e.g., email addresses, preferences) is a basic stub. This means that while the database functionality is simulated or minimally implemented, it serves the purpose of integrating with the communication service for the sake of demonstrating the workflow without needing a full database setup.

- **Email Client**: Similarly, the email client used for sending notifications to fans is assumed to be a simple stub. This approach allowed me to focus on the integration and message flow rather than the complexities of actual email delivery services. The email client's functionality is therefore simulated to show how emails would be sent out without connecting to a real email service provider.

- **Distribution Clients**: The distribution service, responsible for managing the sending out of CDs and vinyls, interacts with distribution clients that are also assumed to be basic stubs. This means the actual process of managing inventory, packaging, and shipping is abstracted away, allowing the focus to remain on the service's role in the larger system.

### Running the Project with Docker Compose

When setting up and running the project using Docker, it's essential to ensure that any changes made to the Dockerfiles are reflected when you start the services. To achieve this, you should use the following Docker Compose command:

```docker-compose up --build```


## Development Thoughts

### On Reusing Code (DRY)

I'm all for not repeating code when it makes sense, so I put the RabbitMQ stuff and the Release model in a shared library. But, you'll notice some duplication between the communication and distribution services. My thought here is that microservices should really stand on their own as much as possible unless it really makes sense to use a shared library. there's a bit of overlap now, but as they grow, they'll likely develop their own unique logic.

### Handling Connection Issues

I added a back-off retry strategy because we can't always count on RabbitMQ being ready to connect. This way, we avoid some potential headaches. I could have used docker-compose to make the services wait, but I wanted to keep things realistic.

### Using Concurrency and Parallelism

I like to keep things simple and only introduce complexity when it's necessary:

1. Sending messages from the API needed to be concurrent so a hiccup in one queue doesn't stop everything.
2. Sending out mass emails in the communication service also needed some parallel processing. We can't afford to go one by one when we're talking about reaching a lot of fans.

I didn't think it was necessary for the distribution service since we're likely dealing with fewer distributors.

## Areas for Improvement

Even with the best effort, there's always room for improvement:

1. **Testing**: I didn't get around to setting up tests for the RabbitMQ integrations due to the tight schedule the unit tests for business logic are very basic too
2. **Robustness**: No health checks or reconnection strategies for RabbitMQ, and no trace IDs for message tracking and no message ids to ensure.

## Closing Thoughts

Given the quick pace, I'm sure there's room for a bit of polish. I'm open to discussing any points that might need a second look. The project brief was a bit vague, so it was tricky to nail down exactly what was needed. I erred on the side of caution with the starter code, only making essential changes.

I did notice the starter code was not handling bad data, if the intention was for this to be fixed it sadly has not been

## Time Spent

Just to be transparent, I spent just over 3 hours on this project. It was a bit of a rush, but I made the most of the time I had.

