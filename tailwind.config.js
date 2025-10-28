/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.{templ,go}",
    "./*.go"
  ],
  theme: {
    extend: {
      fontSize: {
        '2xs': '0.625rem', // 10px
      },
      colors: {
        gray: {
          850: '#1f2937',
        }
      }
    },
  },
  plugins: [],
}