module.exports = {
  mode: "jit",
  content: ["./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      "waves-0": "url('/waves0.svg')",
      "waves-1": "url('/waves1.svg')",
      "waves-2": "url('/waves2.svg')",
      "waves-3": "url('/waves3.svg')",
      "waves-4": "url('/waves4.svg')"
    },
  },
  plugins: [],
}
