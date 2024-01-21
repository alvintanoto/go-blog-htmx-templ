/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./view/*.{templ,html,js}"],
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
      normal: '400',
      medium: '500',
      semibold: '600'
    },
    colors: {
      white: '#FFFFFF',
      grey: '#F2F3F5',
      "grey-darker": '#d9d9d9',
      primary: '#1677ff',
      success: '#52c41a',
      warning: '#faad14',
      error: '#fff2f0',
      "error-border": '#ffccc7'
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

