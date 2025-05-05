/** @type {import('tailwindcss').Config} */
const colors = require('tailwindcss/colors');

delete colors.lightBlue;
delete colors.warmGray;
delete colors.trueGray;
delete colors.coolGray;
delete colors.blueGray;

module.exports = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './partials/**/*.{js,ts,jsx,tsx,mdx}',
    './src/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    fontSize: {
      xs: ['0.75rem', { lineHeight: '1rem' }],
      sm: ['0.875rem', { lineHeight: '1.5rem' }],
      base: ['1rem', { lineHeight: '2rem' }],
      lg: ['1.125rem', { lineHeight: '1.75rem' }],
      xl: ['1.25rem', { lineHeight: '2rem' }],
      '2xl': ['1.5rem', { lineHeight: '2.5rem' }],
      '3xl': ['2rem', { lineHeight: '2.5rem' }],
      '4xl': ['2.5rem', { lineHeight: '3rem' }],
      '5xl': ['3rem', { lineHeight: '3.5rem' }],
      '6xl': ['3.75rem', { lineHeight: '1' }],
      '7xl': ['4.5rem', { lineHeight: '1' }],
      '8xl': ['6rem', { lineHeight: '1' }],
      '9xl': ['8rem', { lineHeight: '1' }],
    },
    // colors: {
    //   ...colors,
    //   transparent: 'transparent',
    //   current: 'currentColor',
    //   'landing-bg': '#12161E',
    //   primary: {
    //     DEFAULT: '#1f2937',
    //     hover: '#374151',
    //     selected: '#111827',
    //   },
    //   secondary: {
    //     DEFAULT: '#4f46e5',
    //   },
    //   txtPrimary: {
    //     DEFAULT: '#d1d5db',
    //     hover: '#ffffff',
    //     selected: '#ffffff',
    //     landing: '#BBCADE',
    //     // disabled: '#06b6d4',
    //   },
    // },
    colors: {
      transparent: 'transparent',
      current: 'currentColor',
      black: colors.black,
      white: colors.white,
      gray: colors.slate,
      primary: colors.cyan,
      green: colors.emerald,
      yellow: colors.yellow,
      red: colors.rose,
      indigo: colors.indigo,
      orange: colors.orange,
    },
    screens: {
      xs: '10px',
      // => @media (min-width: 10px) { ... }

      sm: '640px',
      // => @media (min-width: 640px) { ... }

      md: '768px',
      // => @media (min-width: 768px) { ... }

      lg: '1024px',
      // => @media (min-width: 1024px) { ... }

      xl: '1280px',
      // => @media (min-width: 1280px) { ... }

      '2xl': '1536px',
      // => @media (min-width: 1536px) { ... }
    },
  },
  plugins: [
    require('@tailwindcss/forms')({
      // strategy: 'base',
      strategy: 'class',
    }),
  ],
};
