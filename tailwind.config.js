/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/view/**/*.{templ,html,js}"],
  theme: {
    fontSize: {
      sm: ['12px', '20px'],
      base: ['14px', '22px'],
      xl: ['16px', '24px'],
      '2xl': ['20px', '28px'],
      '3xl': ['24px', '32px'],
      '4xl': ['30px', '38px'],
      '5xl': ['38px', '46px'],
      '6xl': ['46px', '54px'],
      '7xl': ['56px', '64px'],
      '8xl': ['68px', '76px'],
    },
    fontWeight: {
      light: '200',
      normal: '400',
      medium: '500',
      semibold: '600'
    },
    colors: {
      // Functional Color
      black: '#000000',
      white: '#FFFFFF',
      transparent: "transparent",
      primary: '#1677ff',
      success: '#52c41a',
      warning: '#faad14',
      danger: '#ff4d4f',

      // Neutral Color (light)
      "heading-text": "#000000E0",
      "text": "#000000E0",
      "secondary-text": "#000000A6",
      "disabled": "#00000040",
      "default-border": "#D9D9D9FF",
      "separator": "#0505050F",
      "layout-background": "#F5F5F5FF",

      // Neutral Color (dark)
      "dark-heading-text": "#FFFFFFD9",
      "dark-text": "#FFFFFFD9",
      "dark-secondary-text": "#FFFFFFA6",
      "dark-disabled": "#FFFFFF40",
      "dark-default-border": "#424242FF",
      "dark-separator": "#FDFDFD1F",
      "dark-layout-background": "#282C34"

      // Color Theme
    },
    spacing: {
      '1': '8px',
      '2': '12px',
      '3': '20px',
      '4': '32px',
      '5': '48px',
      '6': '80px',
    },
    extend: {},
  },
  plugins: [],
}

