# Instructions for Participants:

Participants receive the names of two images: matixmedia/docker-ctf-main:latest and matixmedia/docker-ctf-data-provider:latest.

-   **Step 1:** Logs & Inspect: Start the main container. The logs reveal that you should use `docker inspect` to find the port (`8989`).
-   **Step 2:** Port Mapping: Restart the main container with `-p 8989:8989`. The website shows an error message that the `data-provider-svc` is not reachable.
-   **Step 3:** Inspect for Label: The hint on the website leads you to use docker inspect again to find the label `ctf.data-provider.host`. The value is `data-provider-svc`.
-   **Step 4:** The Network Problem: Start the second container (`data-provider`) with `--name data-provider-svc`. The connection still fails.
-   **Step 5:** Docker Network: Create a network with `docker network create ctf-net`. Restart both containers and add them to the network with `--network ctf-net`. Now the connection is successful!
-   **Step 6:** Volumes: The website now asks for a password and gives the hint to mount a volume (`-v $(pwd)/secrets:/secrets`) to find it.
-   **Step 7:** Find & Enter Password: After restarting the main container with the volume mount, the file `password.txt` is on the host machine. The participant reads the password (`SUPER_GEHEIM_123`) and enters it on the website.
-   **Step 8:** Final Exec: After entering the correct password, the website displays the final hint: Run `/app/app --show-flag` in the `ctf-main` container. This must be run with an interactive terminal (`-it`).
-   **Goal:** The flag `FLAG{D0CK3R_PR0F1_MIT_FLAG}` is displayed.
