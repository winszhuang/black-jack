/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    fontFamily: {
      sans: ['Poppins', '"Noto Sans TC"', '"Noto Sans"']
    },
    colors: {
      primary: {
        50: '#43506b',
        80: '#46556c',
        100: '#293449',
        200: '#1d2738',
        300: '#151e2d',
        400: '#0e1725'
      },
      white: {
        50: '#ffffff',
        100: '#aeb2b7'
      },
      red: {
        50: '#d9113a',
        100: '#c40a35'
      },
      yellow: {
        50: '#ffd000',
        100: '#e8aa02'
      },
      dark: {
        50: '#F9FAFB',
        100: '#EEEFF7',
        200: '#DCDDED',
        300: '#BCBED1',
        400: '#A1A3BF',
        500: '#6E7191',
        600: '#4B4D6A',
        700: '#37394F',
        800: '#222330',
        900: '#16161F'
      }
    }
  },
  plugins: []
}
