import { LoadingSpinner } from "../molecules/LoadingSpinner";

interface Props {
  showLoading?: boolean;
  hasNextPage?: boolean;
  ref?: React.Ref<HTMLDivElement>;
}

export const BottomOfPage = ({ showLoading = false, hasNextPage = true, ref }: Props) => {
  return (
    <div className={hasNextPage ? "sticky left-0 h-24" : ""} ref={ref}>
      {showLoading && <LoadingSpinner />}
    </div>
  );
}
