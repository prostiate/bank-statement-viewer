import { type Status } from "./CsvUploader";

export default function StatusBadge({ status }: { status: Status }) {
  const getStatusClass = (status: Status): string => {
    switch (status) {
      case "SUCCESS":
        return "success";
      case "PENDING":
        return "warning";
      case "FAILED":
      default:
        return "error";
    }
  };

  const statusClass = getStatusClass(status);

  return <span className={`badge badge-${statusClass}`}>{status}</span>;
}
