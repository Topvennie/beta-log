import { Card, createTheme, Grid, Group, Modal, Notification, rem, SegmentedControl, Select, SimpleGrid, Stack, Title } from "@mantine/core";
import { createElement } from "react";
import { LuSearch } from "react-icons/lu";

const theme = createTheme({
  fontFamily: "Inter, sans-serif",
  fontFamilyMonospace: "monospace",
  defaultRadius: "sm",
  autoContrast: true,

  fontSizes: {
    xs: rem(11),
    sm: rem(12),
    md: rem(13),
    lg: rem(14),
    xl: rem(18),
  },

  lineHeights: {
    xs: rem(14),
    sm: rem(16),
    md: rem(18),
    lg: rem(20),
    xl: rem(24),
  },

  headings: {
    fontFamily: "Inter, sans-serif",
    fontWeight: "600",
    sizes: {
      h1: {
        fontSize: rem(24),
        lineHeight: rem(32),
      },
      h2: {
        fontSize: rem(18),
        lineHeight: rem(24),
      },
      h3: {
        fontSize: rem(14),
        lineHeight: rem(20),
      },
    },
  },

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

  components: {
    Title: Title.extend({
      defaultProps: {
        textWrap: "wrap",
      },
      styles: {
        root: {
          overflowWrap: "anywhere",
          whiteSpace: "pre-wrap",
        },
      },
    }),
    Card: Card.extend({
      defaultProps: {
        shadow: "sm",
        padding: "lg",
      },
    }),
    Modal: Modal.extend({
      defaultProps: {
        centered: true,
        size: "xl",
      },
    }),
    Select: Select.extend({
      defaultProps: {
        searchable: true,
        allowDeselect: true,
        leftSection: createElement(LuSearch),
        className: "grow",
      },
    }),
    SegmentedControl: SegmentedControl.extend({
    }),
    Notification: Notification.extend({
      defaultProps: {
        withBorder: true,
      }
    }),
    Stack: Stack.extend({
      defaultProps: {
        gap: "xs",
      }
    }),
    Group: Group.extend({
      defaultProps: {
        gap: "xs",
      }
    }),
    Grid: Grid.extend({
      defaultProps: {
        gap: "xs",
      }
    }),
    SimpleGrid: SimpleGrid.extend({
      defaultProps: {
        spacing: "xs",
      }
    }),
  }

});

export default theme;
