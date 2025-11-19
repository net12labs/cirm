// Define custom element x-domain that extends div
class XDomain extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        // Element is added to the DOM
        const domain = this.getAttribute('domain');
        if (domain) {
            this.dataset.domain = domain;
        }
    }

    static get observedAttributes() {
        return ['domain'];
    }

    attributeChangedCallback(name, oldValue, newValue) {
        if (name === 'domain' && oldValue !== newValue) {
            this.dataset.domain = newValue;
        }
    }
}

// Register the custom element
customElements.define('x-domain', XDomain);

// Reusable component loader with script execution support
async function loadComponent(containerId, endpoint, componentName, useShadowDOM = false) {
    const container = document.getElementById(containerId);
    try {
        const response = await fetch(endpoint);
        if (response.ok) {
            const html = await response.text();

            if (useShadowDOM) {
                // Create shadow DOM for isolated styles
                const shadowRoot = container.attachShadow({ mode: 'open' });

                // Parse the HTML
                const parser = new DOMParser();
                const doc = parser.parseFromString(html, 'text/html');

                // Append all body children to shadow root
                Array.from(doc.body.childNodes).forEach(node => {
                    shadowRoot.appendChild(node.cloneNode(true));
                });

                // Execute inline scripts in shadow DOM context
                const scripts = shadowRoot.querySelectorAll('script');
                scripts.forEach(oldScript => {
                    const newScript = document.createElement('script');

                    // Copy attributes
                    Array.from(oldScript.attributes).forEach(attr => {
                        newScript.setAttribute(attr.name, attr.value);
                    });

                    // Copy content and modify to use shadowRoot
                    let scriptContent = oldScript.textContent;
                    // Replace document.getElementById with shadowRoot.getElementById
                    scriptContent = scriptContent.replace(/document\.getElementById\(/g, 'shadowRoot.getElementById(');
                    scriptContent = scriptContent.replace(/document\.querySelector\(/g, 'shadowRoot.querySelector(');
                    scriptContent = scriptContent.replace(/document\.querySelectorAll\(/g, 'shadowRoot.querySelectorAll(');
                    // Keep document.addEventListener for document-level events

                    // Wrap in IIFE with shadowRoot reference
                    newScript.textContent = `(function(shadowRoot) { ${scriptContent} })(document.getElementById('${containerId}').shadowRoot);`;

                    // Append to document body to execute
                    document.body.appendChild(newScript);
                    oldScript.remove();
                });
            } else {
                // Regular loading without shadow DOM
                const parser = new DOMParser();
                const doc = parser.parseFromString(html, 'text/html');

                // Clear container
                container.innerHTML = '';

                // Append all body children to container
                Array.from(doc.body.childNodes).forEach(node => {
                    container.appendChild(node.cloneNode(true));
                });

                // Execute inline scripts
                const scripts = container.querySelectorAll('script');
                scripts.forEach(oldScript => {
                    const newScript = document.createElement('script');

                    // Copy attributes
                    Array.from(oldScript.attributes).forEach(attr => {
                        newScript.setAttribute(attr.name, attr.value);
                    });

                    // Copy content
                    newScript.textContent = oldScript.textContent;

                    // Replace old script with new one to execute it
                    oldScript.parentNode.replaceChild(newScript, oldScript);
                });
            }

        } else {
            container.innerHTML = `<p style="color: #d32f2f;">Failed to load ${componentName}.</p>`;
        }
    } catch (error) {
        container.innerHTML = `<p style="color: #d32f2f;">Error loading ${componentName}.</p>`;
        console.error(`Error loading ${componentName}:`, error);
    }
}