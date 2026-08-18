# Image Publication Contract

For repository `OWNER/REPOSITORY`:

- Backend: `ghcr.io/OWNER/REPOSITORY-backend`
- Frontend: `ghcr.io/OWNER/REPOSITORY-frontend`

On default branch publication:

- `main`
- `sha-<short commit>`
- `latest`

Pull Request runs build both images with `push=false` and no registry login.
