// composables/useEventImageManager.ts

import { ref } from "vue";
import { useMutation } from "@vue/apollo-composable";
import { useToast } from "vue-toastification";
import {
  UPLOAD_EVENT_IMAGES,
  DELETE_EVENT_IMAGE,
  SET_FEATURED_IMAGE,
  removeFeaturedMutation,
} from "~/graphql/eventMutations";

export const useEventImageManager = () => {
  const toast = useToast();
  const uploading = ref(false);
  const deleting = ref(false);
  const updatingFeatured = ref(false);
  const uploadProgress = ref(0);
  const errors = ref<string[]>([]);

  const { mutate: uploadMutation } = useMutation(UPLOAD_EVENT_IMAGES);
  const { mutate: deleteMutation } = useMutation(DELETE_EVENT_IMAGE);
  const { mutate: setFeaturedMutation } = useMutation(SET_FEATURED_IMAGE);
  const { mutate: removeFeatured } = useMutation(removeFeaturedMutation);

  // Read raw local files into base64 data channels
  const fileToBase64 = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result as string);
      reader.onerror = (error) => reject(error);
    });
  };

  // removefromfeature
  const uploadImages = async (
    eventId: string | null,
    files: File[],
  ): Promise<string[]> => {
    if (!eventId || files.length === 0) {
      toast.error("No event ID or files provided");
      return [];
    }

    uploading.value = true;
    uploadProgress.value = 0;
    errors.value = [];
    const allUploadedUrls: string[] = [];

    try {
      // Map files array asynchronously into strings
      const imageData = await Promise.all(
        files.map((file) => fileToBase64(file)),
      );

      // Chunks payload requests to avoid running into payload body limitations (10 max)
      const chunkSize = 10;

      for (let i = 0; i < imageData.length; i += chunkSize) {
        const chunk = imageData.slice(i, i + chunkSize);
        uploadProgress.value = Math.round((i / imageData.length) * 100);

        // Dispatches payload structural changes directly to Apollo client instance
        const result = await uploadMutation({
          eventId: eventId,
          images: chunk,
        });

        const response = result?.data?.uploadEventImages;

        if (response?.success) {
          allUploadedUrls.push(...(response.urls || []));
        } else {
          errors.value.push(
            response?.message || "Failed processing chunk sequence upload.",
          );
        }
      }

      uploadProgress.value = 100;

      if (allUploadedUrls.length > 0) {
        toast.success(
          `Successfully uploaded ${allUploadedUrls.length} image(s)`,
        );
      }

      return allUploadedUrls;
    } catch (error: any) {
      console.error("Upload Process Trace Error:", error);
      toast.error(
        error.message ||
          "Failed to complete image files upload processing pipeline.",
      );
      return [];
    } finally {
      uploading.value = false;
    }
  };

  const setFeaturedImage = async (
    eventId: string | null,
    imageId: string,
  ): Promise<any | null> => {
    console.log("from set featured ", eventId, imageId);
    if (!eventId || !imageId) {
      toast.error("Missing event ID or image ID");
      return null;
    }

    updatingFeatured.value = true;

    try {
      console.log("Setting featured image:", { eventId, imageId });

      const result = await setFeaturedMutation({ eventId, imageId });
      console.log("Set featured result:", result);

      const response = result?.data?.setFeaturedImage;
     // Hasura returns data inside the table name operation slot directly
      const affectedRows = result?.data?.update_event_images?.affected_rows;
      const updatedImages = result?.data?.update_event_images?.returning;

      // if (response?.success) {
      if (affectedRows && affectedRows > 0) {

        toast.success("Featured image updated");
        return response?.data?.image || null;
      } else {
        toast.error(response?.message || "Failed to set featured image");
        return null;
      }
    } catch (error: any) {
      console.error("Set featured error:", error);
      toast.error(error.message || "Failed to set featured image");
      return null;
    } finally {
      updatingFeatured.value = false;
    }
  };

  const removeFeaturedImage = async (
    eventId: string | null,
    imageId: string,
  ): Promise<any | null> => {
    if (!eventId || !imageId) {
      toast.error("Missing event ID or image ID");
      return null;
    }

    updatingFeatured.value = true;

    try {
      const result = await removeFeatured({
        eventId,
        imageId,
      });

      // Hasura returns data inside the table name operation slot directly
      const affectedRows = result?.data?.update_event_images?.affected_rows;
      const updatedImages = result?.data?.update_event_images?.returning;

      if (affectedRows && affectedRows > 0) {
        toast.success("Featured image status removed");
        // Return the first modified item from the array
        return updatedImages[0];
      } else {
        toast.error("No record found or updated");
        return null;
      }
    } catch (error: any) {
      console.error("Remove featured error:", error);
      toast.error(error.message || "Failed to remove featured image status");
      return null;
    } finally {
      updatingFeatured.value = false;
    }
  };

  const deleteImage = async (imageId: string): Promise<boolean> => {
    if (!imageId) {
      toast.error("No image ID provided");
      return false;
    }

    deleting.value = true;

    try {
      const result = await deleteMutation({
        imageId,
      });

      // Hasura returns the affected_rows inside delete_event_images
      const affectedRows = result?.data?.delete_event_images?.affected_rows;

      if (affectedRows && affectedRows > 0) {
        // Success! Return true so removeExistingImage can splice it from the array
        return true;
      } else {
        toast.error("Image not found or already deleted");
        return false;
      }
    } catch (error: any) {
      console.error("Delete error:", error);
      toast.error(error.message || "Failed to delete image");
      return false;
    } finally {
      deleting.value = false;
    }
  };
  // Set featured image
  const validateImages = (
    files: File[],
  ): { valid: File[]; errors: string[] } => {
    const valid: File[] = [];
    const errors: string[] = [];

    const maxSize = 10 * 1024 * 1024; // 10MB
    const allowedTypes = ["image/jpeg", "image/png", "image/webp", "image/gif"];

    for (const file of files) {
      if (!allowedTypes.includes(file.type)) {
        errors.push(`${file.name} is not a valid image format`);
        continue;
      }

      if (file.size > maxSize) {
        errors.push(`${file.name} exceeds 10MB limit`);
        continue;
      }

      valid.push(file);
    }

    return { valid, errors };
  };

  return {
    uploading,
    deleting,
    updatingFeatured,
    uploadProgress,
    errors,
    uploadImages,
    deleteImage,
    setFeaturedImage,
    validateImages,
    fileToBase64,
    removeFeaturedImage,
  };
};
