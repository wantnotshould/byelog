# byelog

**Stop tailing, start querying.**

## The "I'm a Clown" Story 🤡

I started this project because I was sick of staring at endless `tail -f` logs in OpenResty. I thought, "Hey, I'll just write a JS SDK and collect everything myself! I'm a genius!"

**I was wrong.** I quickly realized that real threats (bots, crawlers, hackers) don't give a damn about my JavaScript. They don't run it. So for security, this project is basically a toy. OpenResty is still the king there.

## Why this?

After the "clown moment," I figured: "Screw it, I'll just turn this into a learning project to play with some cool tech."

It’s my personal playground for messing with Go, Kafka, and seeing how much data I can collect before I have to care about DB indexing.

## Quick Start

It's just a simple Go server and a tiny JS file.

### 1. Fire up your Kafka broker.

### 2. Run the Go

You can grab the pre-built binary for your OS from [Releases](https://github.com/wantnotshould/byelog/releases), or just build it yourself with Ent and Wire if you prefer.

### 3. Embed the SDK

```html
<script src="https://cdn.jsdelivr.net/gh/wantnotshould/byelog@main/byelog.js"></script>
<script>
    ByeLog.init('your-project-uuid', { 
        serverUrl: 'https://your-log-server.com' 
    });
</script>
```

## ⚖️ License

MIT License. See [LICENSE](./LICENSE).