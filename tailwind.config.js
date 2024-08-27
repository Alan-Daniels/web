/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./internal/**/*.{go,js,templ,html}"
    ],
    theme: {
      extend: {},
        colors: {
            transparent: 'transparent',
            current: 'currentColor',
            'white': '#ffffff',
            'muave': '#cba6f7',
            'base': '#1e1e2e',
            'text': "#cdd6f4",
            'surface': {
                100: "#313244",
                200: "#45475a",
                300: "#585b70",
            },
            'overlay': {
                100: "#6c7086",
                200: "#7f849c",
                300: "#9399b2",
            },

        }
    },
    plugins: [],
  }
