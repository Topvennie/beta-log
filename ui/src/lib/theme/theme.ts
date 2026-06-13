import { createTheme } from "@mantine/core";

const theme = createTheme({
  fontFamily: 'Inter, sans-serif',
  defaultRadius: "md",

  // Same breakpoint values as tailwind
  breakpoints: {
    xs: "36em",
    sm: "40em",
    md: "48em",
    lg: "64em",
    xl: "80em",
    xxl: "96em",
    xxxl: "142em",
    xxxxl: "172em",
  },
});

export default theme
