import { AlertTriangle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { useTranslation } from "react-i18next";

interface DeleteConfirmModalProps {
  isOpen: boolean;
  title: string;
  description: string;
  onConfirm: () => void;
  onCancel: () => void;
  isLoading?: boolean;
}

export function DeleteConfirmModal({
  isOpen,
  title,
  description,
  onConfirm,
  onCancel,
  isLoading = false,
}: DeleteConfirmModalProps) {
  const { t } = useTranslation();

  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm"
        onClick={onCancel}
      />

      {/* Modal */}
      <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
        <Card className="border-red-500/30 bg-card/95 backdrop-blur w-full max-w-sm shadow-2xl">
          <CardContent className="pt-6">
            {/* Icon and Title */}
            <div className="flex gap-3 mb-4">
              <div className="flex-shrink-0">
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-red-500/10">
                  <AlertTriangle className="h-6 w-6 text-red-600 dark:text-red-400" />
                </div>
              </div>
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-red-600 dark:text-red-400">
                  {title}
                </h3>
              </div>
            </div>

            {/* Description */}
            <p className="text-sm text-muted-foreground mb-6 leading-relaxed">
              {description}
            </p>

            {/* Actions */}
            <div className="flex gap-3">
              <Button
                onClick={onCancel}
                variant="outline"
                disabled={isLoading}
                className="flex-1 border-border text-foreground hover:bg-secondary/50"
              >
                {t("common.cancel") || "Cancel"}
              </Button>
              <Button
                onClick={onConfirm}
                disabled={isLoading}
                className="flex-1 bg-red-600 hover:bg-red-700 text-white"
              >
                {isLoading
                  ? t("common.deleting") || "Deleting..."
                  : t("common.delete") || "Delete"}
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </>
  );
}
