<div align="right">
  <img src="https://img.shields.io/badge/-English-F54A00?style=for-the-badge" alt="English" />
  <a title="zh-CN" href="./README.md">  <img src="https://img.shields.io/badge/-%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87-545759?style=for-the-badge" alt="简体中文"></a>
</div>

<div align="center">
  <img alt="Ech0" src="./docs/imgs/logo.svg" width="150">

  [Live Demo](https://memo.vaaat.com/) | [Official Website](https://echo.soopy.cn/) | [Official Docs](https://echodoc.soopy.cn/) | [Ech0 Hub](https://echohub.soopy.cn/)

  # Ech0
</div>

> Open source, self-hosted, lightweight publishing platform focused on the flow of ideas

Ech0 is a lightweight, self-hosted platform designed for quick sharing of your ideas, texts, and links. With an intuitive interface, you can easily manage your content, ensuring complete data control and seamless connection with the world anytime, anywhere.

![Screenshot](./docs/imgs/screenshot_mockup.png)

---

## Core Features

☁️ **Atomic Lightweight**: Uses less than **15MB** RAM, image size under **40MB**, stores data in a single SQLite file  
🚀 **Fast Deployment**: No configuration needed, just one command from install to use  
✍️ **Distraction-Free Writing**: Clean online Markdown editor, **supports rich Markdown plugins and preview**  
📦 **Data Sovereignty**: All content stored locally in SQLite, supports RSS subscription  
🎉 **Free Forever**: Open-sourced under the AGPL-3.0 license, no tracking / subscription / service dependency  
🌍 **Cross-Platform**: Fully responsive across mobile, tablet, and desktop browsers  
👾 **PWA Support**: Can be installed as a web app  
📝 **Built-in Todo Manager**: Easily manage daily tasks and track progress efficiently  
🔗 **Ech0 Connect**: Brand-new content aggregation and interconnection system; supports federation, content subscription and sync across multiple instances  
🎵 **Seamless Music Integration**: Lightweight music player with local audio parsing, immersive background playback and focus mode  
🎥 **Instant Video Sharing**: Native Bilibili/YouTube video parsing  
🃏 **Rich Shortcut Cards**: One-click sharing of rich media like website links, GitHub repos, etc. for vivid display  
⚙️ **Advanced Customization**: Power users can customize styles and scripts for expressive sharing

---

## 3-Second Quick Deployment

<!-- ### 🧙 One-Click Script Deployment (Recommended)
```shell
curl -fsSL "http://echo.soopy.cn/install.sh" -o install_ech0.sh && bash install_ech0.sh
``` -->

### 🐳 Docker Deployment (Recommended)

```shell
docker run -d \
  --name ech0 \
  -p 6277:6277 \
  -v /opt/ech0/data:/app/data \
  -e JWT_SECRET="Hello Echos!" \
  sn0wl1n/ech0:latest
```

> 💡 After deployment, visit `ip:6277` to use
> 🚷 It is recommended to change `"Hello Echos!"` in `JWT_SECRET` for better security
> 📍 The first registered account becomes the administrator (currently only admins can publish content)
> 🎈 Data is stored under `/opt/ech0/data`

### 🐋 Docker Compose Deployment

Create a new directory and put the `docker-compose.yml` file inside it.

Then run the following command in that directory:

```shell
docker-compose up -d
```

---

## How to Update

### 🔄 Update with Docker

```shell
# Stop the current container
docker stop ech0

# Remove the container
docker rm ech0

# Pull the latest image
docker pull sn0wl1n/ech0:latest

# Run the updated container
docker run -d \
  --name ech0 \
  -p 6277:6277 \
  -v /opt/ech0/data:/app/data \
  -e JWT_SECRET="Hello Echos" \
  sn0wl1n/ech0:latest
```

### 💎 Update with Docker Compose

```shell
# Go to the compose file directory
cd /path/to/compose

# Pull latest image and restart
docker-compose pull && \
docker-compose up -d --force-recreate

# Clean old images
docker image prune -f
```

---

# 🦖 Roadmap

* [x] Rewrite frontend using plain Vue 3
* [x] Fix a security issue
* [x] Refactor backend for elegance and efficiency
* [x] Resolve cross-platform issues
* [x] Redesign logo using Figma
* [ ] Polish UI details and add more practical features
* [ ] Performance optimization & UI beautification

---

# ❓ FAQ

1. **What is Ech0?**
   Ech0 is a lightweight open-source self-hosted platform designed for quick idea sharing. It provides a clean interface and distraction-free writing, with all data stored locally for full control.

2. **What is Ech0 not?**
   Ech0 is not a traditional note-taking tool (like Obsidian or Notion). It functions more like a personal feed or “status update” platform.

3. **Is Ech0 free?**
   Yes, Ech0 is completely free and open-source under the AGPL-3.0 License. No ads, no tracking, no subscriptions, and no service dependencies.

4. **How to back up and restore data?**
   All content is stored in a local SQLite file. Simply back up the `/opt/ech0/data` directory. To restore, just copy the backup files back.

5. **Does Ech0 support RSS?**
   Yes, Ech0 supports RSS subscriptions for content updates.

6. **Why do I get a "contact admin" error when publishing?**
   Only administrators can publish content. The first registered user becomes the admin. Assign admin rights to others via the settings.

7. **Why isn’t there a full permission system?**
   Ech0 is designed to be simple and lightweight, avoiding complexity in permission management. Only admin vs non-admin roles exist for now. Simplicity over complexity.

8. **Why can’t others see my Connect avatar?**
   Make sure you set the `Instance URL` under `System Settings` to your deployment domain, e.g., `https://memo.vaaat.com`. It must include `http` or `https`.

9. **What is the MetingAPI setting for?**
   This API is used to fetch streaming music links for the music card feature. If unset, the default hosted version will be used (deployed via Vercel).

10. **Why do some Connect entries not show?**
    The backend fetches data from all connected instances. If one is down or unreachable, it will be excluded from results.

11. **What content is not recommended for Ech0?**
   Ech0 published content consists of three parts: text, images, and extended content (such as music, video player cards). Ech0 does not recommend publishing dense content that simultaneously includes `text + images + extended content`, as this violates some of Ech0's design principles. Additionally, publishing extended content or lengthy articles is not recommended at any time.

12. **How do I enable the comment feature?**  
   Go to the settings page and enter your deployed Twikoo backend address in the `Comment API` field to enable comments automatically. Currently, only [Twikoo](https://twikoo.js.org/) is supported.

---

# 📢 Feedback & Contributions

If you encounter bugs, please report them via [issues](https://github.com/lin-snow/Ech0/issues). For new feature suggestions or improvements, join the discussion in [discussions](https://github.com/lin-snow/Ech0/discussions).

---

# 🪅 Project Architecture

![Architecture Diagram](./docs/imgs/Ech0技术架构图.svg)

> by Excalidraw

---

# 🛠️ Development

## **Backend Requirements:**

📌 **Go 1.24.3+**

📌 **C Compiler**
Required for libraries like `go-sqlite3` that use CGO:

* Windows:

  * [MinGW-w64](https://winlibs.com/)
  * Add `bin` folder to PATH
* macOS: `brew install gcc`
* Linux: `sudo apt install build-essential`

📌 **Google Wire**
Install [wire](https://github.com/google/wire) for dependency injection code generation:

* `go install github.com/google/wire/cmd/wire@latest`

📌 **Golangci-Lint**  
Install [Golangci-Lint](https://golangci-lint.run/) for linting and formatting:  
- Run `golangci-lint run` in the project root directory for linting  
- Run `golangci-lint fmt` in the project root directory for formatting  

## **Frontend Requirements:**

📌  **NodeJS v23.11.1, PNPM v10.13.1**

> Tip: Use [fnm](https://github.com/Schniz/fnm) for managing multiple Node versions.

---

## **Start Frontend & Backend Together**

**Step 1: Backend (in Ech0 root directory):**

```shell
go run cmd/ech0/main.go # Compile and start backend
```

> If DI has changed, regenerate `wire_gen.go` in `ech0/internal/di/` using `wire` command.

**Step 2: Frontend (new terminal):**

```shell
cd web # Enter frontend folder

pnpm install # Run if dependencies are not installed

pnpm dev # Start frontend preview
```

**Step 3: Access:**
Frontend: [http://localhost:5173](http://localhost:5173)
Backend: [http://localhost:6277](http://localhost:6277) (default)

---

# 🥰 Acknowledgements

* Thanks to [Gin](https://github.com/gin-gonic/gin) for a high-performance backend framework
* Thanks to [Md-Editor-V3](https://github.com/imzbf/md-editor-v3) for an amazing Markdown editor
* Thanks to [Figma](https://www.figma.com/) for an easy-to-use design tool
* Thanks to [VSCode](https://code.visualstudio.com/) and [Jetbrain GoLand](https://www.jetbrains.com/) for excellent developer tools
* Thanks to community users for valuable feedback and improvements
* Thanks to my roommate for designing the initial logo
* Thanks to all contributors and supporters from the open-source community

---

# ✨ growth of Star

<a href="https://www.star-history.com/#lin-snow/Ech0&Timeline">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=lin-snow/Ech0&type=Timeline&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=lin-snow/Ech0&type=Timeline" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=lin-snow/Ech0&type=Timeline" />
 </picture>
</a>

---

# ☕ Support

🌟 If you like **Ech0**, feel free to star the project! 🚀

Ech0 is fully open-source and free. Maintaining and improving it requires community support. If this project helped you, consider donating to support its development. Every bit of encouragement fuels our progress!
Donate via the QR code and leave your GitHub name as a note—you'll be acknowledged in the main `README.md`.

|                  Platform                  |                         QR Code                        |
| :----------------------------------------: | :----------------------------------------------------: |
| [**Afdian**](https://afdian.com/a/l1nsn0w) | <img src="./docs/imgs/pay.jpeg" alt="Pay" width="200"> |

---

```golang
███████╗     ██████╗    ██╗  ██╗     ██████╗ 
██╔════╝    ██╔════╝    ██║  ██║    ██╔═████╗
█████╗      ██║         ███████║    ██║██╔██║
██╔══╝      ██║         ██╔══██║    ████╔╝██║
███████╗    ╚██████╗    ██║  ██║    ╚██████╔╝
╚══════╝     ╚═════╝    ╚═╝  ╚═╝     ╚═════╝ 
``` 