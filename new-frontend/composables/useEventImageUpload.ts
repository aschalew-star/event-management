// composables/useEventImageUpload.ts

import { ref, computed } from 'vue';
import { useMutation } from '@vue/apollo-composable';
import { useToast } from 'vue-toastification';
import { UPLOAD_EVENT_IMAGES, DELETE_EVENT_IMAGE, SET_FEATURED_IMAGE } from '~/graphql/eventMutations';

export const useEventImageUpload = () => {
  const toast = useToast();
  const uploading = ref(false);
  const uploadProgress = ref(0);
  const errors = ref<string[]>([]);

  const { mutate: uploadMutation } = useMutation(UPLOAD_EVENT_IMAGES);
  const { mutate: deleteMutation } = useMutation(DELETE_EVENT_IMAGE);
  const { mutate: setFeaturedMutation } = useMutation(SET_FEATURED_IMAGE);

  // Convert file to base64
  const fileToBase64 = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result as string);
      reader.onerror = (error) => reject(error);
    });
  };

  // Upload images
  const uploadImages = async (eventId: string, files: File[]) => {
    if (!eventId || files.length === 0) {
      toast.error('No event ID or files provided');
      return null;
    }

    uploading.value = true;
    uploadProgress.value = 0;
    errors.value = [];

    try {
      // Convert files to base64
      const fileData = await Promise.all(
        files.map(async (file) => ({
          file: await fileToBase64(file),
          filename: file.name
        }))
      );

      // Upload via Hasura action
      const result = await uploadMutation({
        eventId,
        files: fileData
      });

      const response = result?.data?.uploadEventImages;

      if (response?.success) {
        toast.success(response.message || 'Images uploaded successfully!');
        
        // Show any partial errors
        if (response.data?.errors?.length > 0) {
          response.data.errors.forEach((err: any) => {
            toast.warning(`Failed to upload ${err.filename}: ${err.error}`);
          });
        }

        return response.data?.images || [];
      } else {
        toast.error(response?.message || 'Failed to upload images');
        return null;
      }

    } catch (error: any) {
      console.error('Upload error:', error);
      toast.error(error.message || 'Failed to upload images');
      return null;
    } finally {
      uploading.value = false;
      uploadProgress.value = 100;
    }
  };

  // Delete image
  const deleteImage = async (imageId: string) => {
    if (!imageId) {
      toast.error('No image ID provided');
      return false;
    }

    try {
      const result = await deleteMutation({ imageId });
      const deleted = result?.data?.delete_event_images_by_pk;

      if (deleted) {
        toast.success('Image deleted successfully');
        return true;
      } else {
        toast.error('Failed to delete image');
        return false;
      }
    } catch (error: any) {
      console.error('Delete error:', error);
      toast.error(error.message || 'Failed to delete image');
      return false;
    }
  };

  // Set featured image
  const setFeaturedImage = async (eventId: string, imageId: string) => {
    if (!eventId || !imageId) {
      toast.error('Missing event ID or image ID');
      return null;
    }

    try {
      const result = await setFeaturedMutation({
        eventId,
        imageId
      });

      const featured = result?.data?.update_event_images_by_pk;

      if (featured) {
        toast.success('Featured image updated');
        return featured;
      } else {
        toast.error('Failed to set featured image');
        return null;
      }
    } catch (error: any) {
      console.error('Set featured error:', error);
      toast.error(error.message || 'Failed to set featured image');
      return null;
    }
  };

  return {
    uploading,
    uploadProgress,
    errors,
    uploadImages,
    deleteImage,
    setFeaturedImage
  };
};