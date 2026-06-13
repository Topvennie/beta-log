import { Button, ButtonProps } from "@mantine/core";
import {
  createLink,
  LinkComponent,
} from "@tanstack/react-router";
import React from "react";

type LinkButtonProps = Omit<ButtonProps, "href">;

const SimpleButtonLinkComponent = React.forwardRef<HTMLAnchorElement, LinkButtonProps>((props, ref) => {
  return <Button ref={ref} component="a" {...props} />;
});

const SimpleLink = createLink(SimpleButtonLinkComponent);

export const LinkButton: LinkComponent<typeof SimpleButtonLinkComponent> = (props) => {
  return <SimpleLink preload="intent" {...props} />;
};
