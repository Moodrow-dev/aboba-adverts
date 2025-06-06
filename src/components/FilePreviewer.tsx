import React, { useImperativeHandle, useState, forwardRef } from "react";
import { Image, message } from "antd";
import { UploadFile } from "antd/lib/upload/interface";

// Экспортируем интерфейс
export interface FileReviewerRef {
    handlePreview: (file: UploadFile) => void;
}

interface FilePreviewerProps {
    onDownloadFile?: (fileId: string) => Promise<Blob>;
}

const FilePreviewer = forwardRef<FileReviewerRef, FilePreviewerProps>(
    ({ onDownloadFile }, ref) => {
        const [previewImage, setPreviewImage] = useState<string>("");
        const [previewOpen, setPreviewOpen] = useState(false);

        const handlePreview = async (file: UploadFile) => {
            if (file.preview) {
                setPreviewImage(file.preview);
                setPreviewOpen(true);
                return;
            }

            if (file.response?.id && onDownloadFile) {
                try {
                    const blob = await onDownloadFile(file.response.id);
                    const previewUrl = await createPreviewFromBlob(blob);
                    setPreviewImage(previewUrl);
                    setPreviewOpen(true);
                } catch (error) {
                    console.error("File loading error:", error);
                    message.error("Ошибка при загрузке файла");
                }
            }
        };

        const createPreviewFromBlob = (blob: Blob): Promise<string> => {
            return new Promise((resolve, reject) => {
                const reader = new FileReader();
                reader.onload = () => resolve(reader.result as string);
                reader.onerror = reject;
                reader.readAsDataURL(blob);
            });
        };

        useImperativeHandle(ref, () => ({
            handlePreview,
        }));

        return (
            <Image
                preview={{
                    visible: previewOpen,
                    onVisibleChange: (visible) => {
                        setPreviewOpen(visible);
                        if (!visible) setPreviewImage("");
                    },
                }}
                src={previewImage}
                style={{ display: "none" }}
            />
        );
    }
);

export default FilePreviewer;