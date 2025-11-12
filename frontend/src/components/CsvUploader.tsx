import React, { useRef, useState } from "react";
import GetBalance from "./GetBalance";
import { useQueryClient } from "@tanstack/react-query";
import GetIssues from "./GetIssues";

export type Status = "SUCCESS" | "PENDING" | "FAILED";
export type TransactionType = "CREDIT" | "DEBIT";
export type Transaction = {
  id: number;
  timestamp: string;
  name: string;
  type: TransactionType;
  amount: number;
  status: Status;
  description: string;
};

function CsvUploader() {
  const queryClient = useQueryClient();
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploadStatus, setUploadStatus] = useState("");
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedFile(event.target.files?.[0] || null);
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      setUploadStatus("Please select a file first.");
      return;
    }

    if (selectedFile.type !== "text/csv") {
      setUploadStatus("Please select a CSV file.");
      return;
    }

    const formData = new FormData();
    formData.append("file", selectedFile);

    try {
      setUploadStatus("Uploading...");

      const response = await fetch("/api/upload", {
        method: "POST",
        body: formData,
      });

      if (response.ok) {
        setSelectedFile(null);
        setUploadStatus("File uploaded successfully!");

        if (fileInputRef.current) {
          fileInputRef.current.value = "";
        }

        const data = await response.json();
        console.log("data:", data);

        await queryClient.invalidateQueries({ queryKey: ["balance"] });
        await queryClient.invalidateQueries({ queryKey: ["issues"] });
      } else {
        const errorData = await response.json();
        console.error("errorData:", errorData);
        setUploadStatus(
          `Upload failed: ${errorData.error || response.statusText}`
        );
      }
    } catch (error) {
      console.error("Uploading error:", error);
      setUploadStatus("An error occurred during upload.");
    }
  };

  return (
    <div>
      <input
        ref={fileInputRef}
        type="file"
        accept=".csv"
        onChange={handleFileChange}
      />
      <button onClick={handleUpload} disabled={!selectedFile}>
        {selectedFile ? "Upload CSV" : "Select a file"}
      </button>
      {uploadStatus && <p>{uploadStatus}</p>}
      <GetBalance />
      <GetIssues />
    </div>
  );
}

export default CsvUploader;
