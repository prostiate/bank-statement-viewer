import { useQuery } from "@tanstack/react-query";
import { type Transaction } from "./CsvUploader";
import StatusBadge from "./StatusBadge";
import { formatNumber } from "../utils/formatNumber";
import { useState } from "react";

type IssuesData = {
  issues: Transaction[];
};

type IssuesResponse = {
  data: IssuesData;
};

type Pagination = {
  page: number;
  rowsPerPage: number;
  sort: keyof Transaction;
  order: "asc" | "desc";
};

function GetIssues() {
  const [pagination, setPagination] = useState<Pagination>({
    page: 0,
    rowsPerPage: 10,
    sort: "timestamp",
    order: "desc",
  });
  const { data, isPending, isError, error } = useQuery<IssuesResponse>({
    queryKey: ["issues"],
    queryFn: () => fetch("/api/issues").then((r) => r.json()),
  });

  const clientSideData = data?.data.issues
    ?.map((issue) => issue)
    .sort((a, b) => {
      if (
        a[pagination.sort as keyof Transaction] <
        b[pagination.sort as keyof Transaction]
      )
        return pagination.order === "asc" ? -1 : 1;
      if (
        a[pagination.sort as keyof Transaction] >
        b[pagination.sort as keyof Transaction]
      )
        return pagination.order === "asc" ? 1 : -1;
      return 0;
    })
    .slice(
      pagination.page * pagination.rowsPerPage,
      (pagination.page + 1) * pagination.rowsPerPage
    );

  if (isPending) return <div>Loading issues...</div>;
  if (isError) return <div>Error: {error?.message}</div>;

  return (
    <div>
      <h2>Issues</h2>
      <table className="table" border={1}>
        <thead>
          <tr>
            <th
              onClick={() =>
                setPagination({
                  ...pagination,
                  sort: "timestamp",
                  order: pagination.order === "asc" ? "desc" : "asc",
                })
              }
            >
              Timestamp{" "}
              <span>
                {pagination.sort === "timestamp"
                  ? pagination.order === "asc"
                    ? "↑"
                    : "↓"
                  : ""}
              </span>
            </th>
            <th
              onClick={() =>
                setPagination({
                  ...pagination,
                  sort: "name",
                  order: pagination.order === "asc" ? "desc" : "asc",
                })
              }
            >
              Name{" "}
              <span>
                {pagination.sort === "name"
                  ? pagination.order === "asc"
                    ? "↑"
                    : "↓"
                  : ""}
              </span>
            </th>
            <th
              onClick={() =>
                setPagination({
                  ...pagination,
                  sort: "type",
                  order: pagination.order === "asc" ? "desc" : "asc",
                })
              }
            >
              Type{" "}
              <span>
                {pagination.sort === "type"
                  ? pagination.order === "asc"
                    ? "↑"
                    : "↓"
                  : ""}
              </span>
            </th>
            <th
              onClick={() =>
                setPagination({
                  ...pagination,
                  sort: "amount",
                  order: pagination.order === "asc" ? "desc" : "asc",
                })
              }
            >
              Amount{" "}
              <span>
                {pagination.sort === "amount"
                  ? pagination.order === "asc"
                    ? "↑"
                    : "↓"
                  : ""}
              </span>
            </th>
            <th
              onClick={() =>
                setPagination({
                  ...pagination,
                  sort: "status",
                  order: pagination.order === "asc" ? "desc" : "asc",
                })
              }
            >
              Status{" "}
              <span>
                {pagination.sort === "status"
                  ? pagination.order === "asc"
                    ? "↑"
                    : "↓"
                  : ""}
              </span>
            </th>
            <th
              onClick={() =>
                setPagination({
                  ...pagination,
                  sort: "description",
                  order: pagination.order === "asc" ? "desc" : "asc",
                })
              }
            >
              Description{" "}
              <span>
                {pagination.sort === "description"
                  ? pagination.order === "asc"
                    ? "↑"
                    : "↓"
                  : ""}
              </span>
            </th>
          </tr>
        </thead>
        <tbody>
          {clientSideData?.length === 0 && (
            <tr>
              <td colSpan={6}>No data found</td>
            </tr>
          )}
          {clientSideData?.map((issue, index) => (
            <tr key={`${index}-${issue.timestamp}`}>
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
      <ul className="pagination">
        {Array.from(
          {
            length: Math.ceil(
              data?.data.issues?.length / pagination.rowsPerPage
            ),
          },
          (_, index) => (
            <li key={index}>
              <button
                onClick={() => setPagination({ ...pagination, page: index })}
                disabled={pagination.page === index}
              >
                {index + 1}
              </button>
            </li>
          )
        )}
      </ul>
    </div>
  );
}

export default GetIssues;
