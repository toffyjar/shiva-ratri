# AstroMonk Mobile App

Expo/React Native client for the Jyotish API. See [AGENTS.md](AGENTS.md) for the Expo SDK version in use — check those docs before writing code.

## Prerequisites

- Node.js and npm
- Docker and Docker Compose
- [Expo Go](https://expo.dev/go) installed on your phone (Android/iOS), on the same Wi-Fi network as your dev machine

## Running on a device

1. **Start the backend services** (from the repo root):

   ```bash
   docker-compose up -d
   ```

   This brings up the Jyotish API (`:9393`), the app backend (`:5656`), and MongoDB. Make sure `api/.env` exists first (copy it from `api/.env.example` if needed).

2. **Point the app at your machine's IP.**

   The phone connects to your dev machine over the LAN, so `localhost` won't work — find your machine's local IP (e.g. `ifconfig` / `ip addr` on Linux/Mac, `ipconfig` on Windows, usually something like `192.168.x.x`) and set it via env vars, or export before starting Expo:

   ```bash
   export EXPO_PUBLIC_JYOTISH_URL=http://<your-lan-ip>:9393
   export EXPO_PUBLIC_APP_BASE_URL=http://<your-lan-ip>:5656/api/v1
   ```

   Note: [pages/chart.tsx](pages/chart.tsx) currently has these URLs hardcoded (`192.168.1.10`) instead of reading the env vars above — update them to your machine's IP directly if you hit connection errors on the chart screen.

3. **Install dependencies and start Expo** (from `mobile-app/`):

   ```bash
   npm install
   npm start
   ```

4. **Scan the QR code** shown in the terminal/browser using the Expo Go app (Android: in-app scanner; iOS: Camera app) to open the app on your device.

## Other useful scripts

```bash
npm run android   # start with Android target
npm run ios       # start with iOS target
npm run web       # start web target
npm test          # run Jest unit tests
```
