/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: ['class', '[data-theme="dark"]'],
  theme: {
    extend: {
      colors: {
        ig: {
          bg: 'var(--ig-bg)',
          card: 'var(--ig-card-bg)',
          border: 'var(--ig-border)',
          text: 'var(--ig-text)',
          'secondary-text': 'var(--ig-secondary-text)',
          blue: 'var(--ig-blue)',
          red: 'var(--ig-red)',
        },
        header: 'var(--header-bg)',
        'card-header': 'var(--card-header-bg)',
        editor: 'var(--editor-bg)',
        output: 'var(--output-text)',
      },
      backgroundImage: {
        'ig-gradient': 'var(--ig-gradient)',
      },
    },
  },
  plugins: [],
}
