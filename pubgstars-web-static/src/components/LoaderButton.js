import React from "react";
import { Button } from "react-bootstrap";
import "./LoaderButton.css";
import Spinner from "react-bootstrap/Spinner";

export default ({
  isLoading,
  text,
  loadingText,
  className = "",
  disabled = false,
  ...props
}) =>
  <Button variant="success" disabled={disabled || isLoading} {...props}>
    {isLoading && <Spinner
                    as="span"
                    animation="border"
                    size="sm"
                    role="status"
                    aria-hidden="true"
                  />}
    {!isLoading ? text :  <span>&nbsp;&nbsp;{loadingText}</span>}
  </Button>
