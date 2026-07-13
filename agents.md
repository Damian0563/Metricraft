### Context of the project

The philosophy behind the project and its intended deployment can be found in the [README.md](https://github.com/damian-bingel/metricraft/blob/main/README.md) file.

There are three main components that constitute the architecture of Metricraft:

1. The **metricraft** directory contains nuxt code for the frontend with tailwind and vue3 components.
2. The **backend** directory contains the golang backend code for the frontend and the API acts as a proxy to the **worker** via grpc.
3. The **worker** directory is golang-powered brains behind processing the actual collected data and organization wide settings.

Besides the three main components, there are also a few other components that are used to support the project:

1. local (hidden in the container) postgresql instance for storing the data related to logs and settings DATABASE_LOGS
2. local (hidden in the container) redis instance for storing the data CACHE and handling sessions
3. external postgreql supabase instance for storing the data related to organizations DATABASE_USERS

### Rules & guidelines

- always aim for modularity and readability of written code
- when defining types define them in a dedicated file, each main component has a coresponding types file
- when designing the frontend code use the color motives used within the project and omit common AI quirks (e.g. glowing html elements, questionable on-hover effects and more), make sure the designed elements are responsive and handle well varying screen sizes
- when building html elements on the frontend optimize for performance and accessibility, in particular optimize for SEO from the start

---
make no mistakes
