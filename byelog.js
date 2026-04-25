!(function (e) {
	const t = {
		config: {serverUrl: '', appId: ''},
		send(e = {}) {
			if (!this.config.appId || !this.config.serverUrl) return
			const n = JSON.stringify({title: e.title || document.title, ...e})
			fetch(this.config.serverUrl, {
				method: 'POST',
				body: n,
				keepalive: !0,
				headers: {'Content-Type': 'application/json', 'Bye-Log-App-Id': this.config.appId}
			}).catch(() => {})
		},
		init(e, n = {}) {
			;((this.config.appId = e),
				(this.config = {...this.config, ...n}),
				'complete' === document.readyState
					? this.send()
					: window.addEventListener('load', () => this.send()))
		}
	}
	window.ByeLog = t
})(window)
