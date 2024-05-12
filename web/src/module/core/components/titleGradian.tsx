import classnames from "classnames";
import React from "react";

interface ITitleGradianProps {
  title: string;
  className?: string;
}

export const TitleGradian: React.FC<ITitleGradianProps> = ({
  title,
  className,
}) => {
  return (
    <div
      className={classnames(
        "bg-gradient-to-r from-[#144EE3] via-[#EB568E] to-[#144EE3] inline-block text-transparent bg-clip-text",
        className,
      )}
    >
      {title}
    </div>
  );
};
