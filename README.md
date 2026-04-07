🕰️ CHRONOS VR | Spatial Time Machine
Chronos VR is an immersive, AI-powered WebXR application that allows users to travel through time. Built with A-Frame for spatial rendering and Go for the backend, the system utilizes Google Gemini to dynamically generate historical contexts, coordinates, and artifact data for any user-requested era.

If a user initiates a "Custom Temporal Jump" (e.g., "Paris 1900"), the Go server acts as an intelligent agent, querying Gemini for structured historical data and simultaneously scraping the Wikipedia API to project a real historical photo onto a holographic screen inside the VR environment.

✨ Features
AI-Generated Environments: Input any historical event or location, and Gemini constructs a complete UI dataset (Location, Coordinates, Artifact Analysis, Mood Colors) on the fly.

Holographic Field References: The Go backend bypasses Wikipedia's strict bot-checkers to dynamically fetch and display relevant historical images on a transparent 3D HUD.

Immersive WebXR: Fully optimized for both standard PC browsers (Mouse Gaze) and native VR Headsets (Meta Quest 3 Laser Controls).

Aggressive Warp Mechanics: Smooth spatial UI animations that fade the HUD and physically move the user's camera to simulate "diving" into history.

🛠️ Tech Stack
Frontend: HTML5, Vanilla JavaScript, A-Frame 1.5.0 (WebXR)

Backend: Go (Golang), net/http router, rs/cors

AI Integration: Google Generative AI SDK (gemini-2.5-flash)

External APIs: Wikimedia Action API

🚀 Getting Started
Prerequisites
Install Go (1.20 or higher).

Get a free API key from Google AI Studio.

Install the "Live Server" extension in VS Code.

Installation & Setup
Clone the repository
git clone https://github.com/yourusername/chronos-vr.git
cd chronos-vr

Configure Environment Variables
Create a .env file in the root directory and add your Gemini API key:
GEMINI_API_KEY=your_actual_api_key_here
PORT=8080

Install Go Dependencies
go mod tidy

Start the Chronos Backend
go run .

Start the Frontend
Right-click index.html in VS Code and select "Open with Live Server". Ensure it opens on http://127.0.0.1:5500 or http://localhost:5500 to satisfy the CORS policy.

🎮 How to Use
On PC/Mac: Click and drag the background to look around. The mouse acts as your pointer. Click the UI panels to switch eras.

On Meta Quest 3: Enter the site using the Meta Quest Browser. Click the VR icon in the bottom right corner to enter immersive mode. Use your hand controllers to point and click.

Custom Jump: Click "Custom Temporal Jump" and enter a location and year (e.g., "The Moon Landing 1969" or "Construction of the Eiffel Tower 1887").

🔒 Security Note
The .env file is explicitly ignored in the .gitignore to prevent API key leaks. Do not commit your Gemini API key to version control.