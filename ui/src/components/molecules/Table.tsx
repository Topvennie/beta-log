import { DataTable, DataTableProps } from "mantine-datatable";
import { LoadingSpinner } from "./LoadingSpinner";

type Props<T> = DataTableProps<T> & Record<string, unknown>;

export const Table = <T,>({ ...props }: Props<T>) => {
  return (
    <DataTable
      customLoader={<LoadingSpinner />}
      minHeight={180}
      backgroundColor="background.0"
      withTableBorder={false}
      textSelectionDisabled={true}
      styles={{
        root: (theme) => ({
          borderRadius: theme.radius.sm,
        }),
        // header: (theme) => ({
        //   background: theme.colors.secondary[1],
        // }),
      }}
      {...props}
    />
  );
};
