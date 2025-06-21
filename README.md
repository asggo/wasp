# WASP
WASP is designed to help you get started building Golang based web applications quickly. I wrote it out a personal need to have a simple base from which to build applications. Every time I thought about building a new webapp, I knew I would have to think about authentication, authorization, and data storage before I could even get to designing the application I wanted to build. WASP has all of that built in and is ready to be extended into whatever application you want to build.

## Getting Started
The first thing you need to do is download the latest release and unzip it in your source code repository. From there you will add handlers to either the `siteRouter`, `adminRouter`, or `userRouter` in `router.go`. The existing handlers live in the `handler` directory and can be modified as needed. In addition, new handlers should be added in the `handler` directory and then called in the appropriate router. Handlers for authenticated endpoints should be added to the `siteRouter` and handlers for administrative endpoints should be added to the `adminRouter`. The application is already built with the necessary authentication, authorization, and session management needed to ensure content in those handlers are protected appropriately.

## Storage
WASP uses the bbolt key value store as its primary storage, but can be extended to use a traditional database as well. If your web application needs new objects such as `posts` or `comments`, they should be added to the `store` directory. If you've never worked with a key value database, I would suggest you give it a try, it is simple, lightweight, and scalable.

## Testing
WASP has tests for all of the core functionality. If you make changes to the core of the application, run the tests to ensure everything works as it should. When you add new objects to the store, I would suggest creating appropriate tests in the `store` directory. You should be able to follow the pattern in the existing tests. If you add new endpoints and handlers, create additional tests following the pattern in the `app_test.go` file and the text files in the `tests` folder.
