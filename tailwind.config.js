/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./components/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [],
  safelist: [
    'bg-gray-100',
    'bg-red-200',
    'bg-yellow-300',
    'bg-green-400',
    'bg-blue-500',
    'bg-indigo-600',
    'bg-purple-700',
    'bg-pink-800',
  ]
}

