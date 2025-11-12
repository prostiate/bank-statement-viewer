import { useQuery } from "@tanstack/react-query";
import { type Transaction } from "./CsvUploader";
import StatusBadge from "./StatusBadge";
import { formatNumber } from "../utils/formatNumber";

type IssuesData = {
  issues: Transaction[];
};

type IssuesResponse = {
  data: IssuesData;
};

function GetIssues() {
  const { data, isPending, isError, error } = useQuery<IssuesResponse>({
    queryKey: ["issues"],
    queryFn: () => fetch("/api/issues").then((r) => r.json()),
  });

  if (isPending) return <div>Loading issues...</div>;
  if (isError) return <div>Error: {error.message}</div>;

  return (
    <div>
      <h2>Issues</h2>
      <table className="table" border={1}>
        <thead>
          <tr>
            <th>Timestamp</th>
            <th>Name</th>
            <th>Type</th>
            <th>Amount</th>
            <th>Status</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          {data.data.issues?.length === 0 && (
            <tr>
              <td colSpan={6}>No data found</td>
            </tr>
          )}
          {data.data.issues?.map((issue) => (
            <tr key={issue.timestamp}>
              <td>{new Date(issue.timestamp).toLocaleString()}</td>
              <td>{issue.name}</td>
              <td>{issue.type}</td>
              <td>{formatNumber(issue.amount)}</td>
              <td>
                <StatusBadge status={issue.status} />
              </td>
              <td>{issue.description}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default GetIssues;
