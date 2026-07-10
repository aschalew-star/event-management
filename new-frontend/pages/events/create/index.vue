<template>
  <div class="container mx-auto max-w-4xl py-8 px-4">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        Create New Event
      </h1>
      <p class="text-gray-600 dark:text-gray-400 mt-2">
        Fill in the details below to create your event
      </p>
    </div>

    <!-- Loading/Submitting State -->
    <div
      v-if="isSubmittingLocal"
      class="mb-6 p-4 bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300 rounded-lg text-sm flex items-center gap-3 border border-blue-200 dark:border-blue-800"
    >
      <i class="fas fa-spinner fa-spin text-lg"></i>
      <span>Creating your event and uploading images... Please wait.</span>
    </div>

    <!-- Error Display -->
    <div
      v-if="error"
      class="mb-6 p-4 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg text-sm border border-red-200 dark:border-red-800 flex items-start gap-2"
    >
      <i class="fas fa-exclamation-circle mt-0.5"></i>
      <span>{{ error }}</span>
    </div>

    <VeeForm
      v-slot="{ errors, isSubmitting }"
      @submit="handleSubmit"
      class="space-y-6"
      :validation-schema="schema"
    >
      <!-- Event Title -->
      <div>
        <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">
          Event Title *
        </label>
        <VeeField
          name="title"
          type="text"
          class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-4 py-2.5 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
          placeholder="Enter your event title"
          v-model="form.title"
        />
        <VeeErrorMessage name="title" class="text-red-500 text-xs mt-1 block" />
      </div>

      <!-- Description -->
      <div>
        <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">
          Description *
        </label>
        <VeeField
          name="description"
          as="textarea"
          rows="5"
          class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-4 py-2.5 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
          placeholder="Describe your event in detail..."
          v-model="form.description"
        />
        <VeeErrorMessage name="description" class="text-red-500 text-xs mt-1 block" />
      </div>

      <!-- Category -->
      <BaseCategorySelect v-model="form.category_id" />

      <!-- Price & Venue -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">
            Price Rules
          </label>
          <div class="flex items-center gap-2 mb-2">
            <input
              id="isFreeCheck"
              type="checkbox"
              class="w-4 h-4 text-indigo-600 rounded border-gray-300 focus:ring-indigo-500"
              v-model="form.is_free"
            />
            <label for="isFreeCheck" class="text-sm text-gray-600 dark:text-gray-400">
              Free Entry Event
            </label>
          </div>
          <VeeField
            v-if="!form.is_free"
            name="price"
            type="number"
            step="0.01"
            min="0"
            class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-4 py-2 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
            placeholder="Price Amount"
            v-model="form.price"
          />
          <VeeErrorMessage name="price" class="text-red-500 text-xs mt-1 block" />
        </div>

        <div>
          <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">
            Venue Name *
          </label>
          <VeeField
            name="venue"
            type="text"
            class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-4 py-2.5 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
            placeholder="Enter venue name"
            v-model="form.venue"
          />
          <VeeErrorMessage name="venue" class="text-red-500 text-xs mt-1 block" />
        </div>
      </div>

      <!-- Location -->
      <div class="p-4 bg-gray-50 dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-700">
        <label class="block text-sm font-bold text-gray-800 dark:text-gray-200 mb-1">
          Location *
        </label>
        <VeeField
          name="address"
          type="text"
          class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-4 py-2.5 text-sm text-gray-900 dark:text-white mb-3 focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
          placeholder="Full address"
          v-model="form.address"
        />
        <VeeErrorMessage name="address" class="text-red-500 text-xs mt-1 block" />

        <EventMapSingle
          :latitude="form.latitude"
          :longitude="form.longitude"
          mode="select"
          @update:location="updateLocation"
        />
      </div>

      <!-- Date & Time -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">
            Event Date *
          </label>
          <VeeField
            name="event_date"
            type="date"
            class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
            v-model="form.event_date"
          />
          <VeeErrorMessage name="event_date" class="text-red-500 text-xs mt-1 block" />
        </div>
        <div>
          <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">
            Start Time
          </label>
          <VeeField
            name="start_time"
            type="time"
            class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
            v-model="form.start_time"
          />
        </div>
        <div>
          <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">
            End Time
          </label>
          <VeeField
            name="end_time"
            type="time"
            class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
            v-model="form.end_time"
          />
        </div>
      </div>

      <!-- Tags -->
      <div>
        <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">
          Tags
        </label>
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-300 dark:border-gray-700 px-1 py-0.5 focus-within:ring-2 focus-within:ring-indigo-500 transition-colors">
          <input
            v-model="tagInput"
            @keydown.enter.prevent="addTag"
            placeholder="Type a tag and press Enter..."
            class="w-full px-3 py-2 text-sm bg-transparent border-none outline-none text-gray-900 dark:text-white"
          />
        </div>
        <div class="flex flex-wrap gap-1.5 mt-2">
          <span
            v-for="(tag, i) in selectedTags"
            :key="i"
            class="inline-flex items-center gap-1 text-xs px-2.5 py-1 rounded-full bg-indigo-50 dark:bg-indigo-900/40 text-indigo-600 dark:text-indigo-400 font-medium"
          >
            {{ tag }}
            <button
              type="button"
              @click="removeTag(i)"
              class="hover:text-red-500 text-[10px] transition-colors"
            >
              <i class="fas fa-times"></i>
            </button>
          </span>
        </div>
      </div>

      <!-- Image Upload Section -->
      <div class="space-y-4">
        <div>
          <div class="flex items-center justify-between mb-2">
            <div>
              <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300">
                Event Images *
              </label>
              <p class="text-xs text-gray-500 dark:text-gray-400">
                Upload up to 5 images. The first image will be featured.
              </p>
            </div>
            <div class="text-sm font-medium" :class="imagePreviews.length > 0 ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-400'">
              {{ imagePreviews.length }}/5
            </div>
          </div>
          
          <!-- Image Upload Area -->
          <div
            class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-4 hover:border-indigo-500 transition-colors"
            :class="imagePreviews.length >= 5 ? 'bg-gray-50 dark:bg-gray-800/50 cursor-not-allowed' : 'cursor-pointer'"
            @dragover.prevent="handleDragOver"
            @drop.prevent="handleDrop"
            @click="handleAreaClick"
          >
            <!-- Images Grid -->
            <div v-if="imagePreviews.length > 0" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
              <div
                v-for="(preview, index) in imagePreviews"
                :key="index"
                class="relative group"
              >
                <div class="relative aspect-square rounded-lg overflow-hidden bg-gray-100 dark:bg-gray-700">
                  <img
                    :src="preview"
                    :alt="`Image ${index + 1}`"
                    class="w-full h-full object-cover"
                  />
                  <!-- Image Number Badge -->
                  <div
                    class="absolute top-1 left-1 bg-black/60 text-white text-[10px] font-bold px-2 py-0.5 rounded-full"
                  >
                    #{{ index + 1 }}
                  </div>
                  <!-- Featured Badge -->
                  <div
                    v-if="index === 0"
                    class="absolute top-8 left-1 bg-indigo-600 text-white text-[10px] font-bold px-2 py-0.5 rounded-full"
                  >
                    Featured
                  </div>
                  <!-- Remove Button -->
                  <button
                    type="button"
                    @click.stop="removeImage(index)"
                    class="absolute top-1 right-1 p-1.5 bg-red-500 text-white rounded-full hover:bg-red-600 transition-colors opacity-0 group-hover:opacity-100"
                  >
                    <i class="fas fa-times text-xs"></i>
                  </button>
                </div>
              </div>
              
              <!-- Add More Button - Only show if less than 5 images -->
              <div
                v-if="imagePreviews.length < 5"
                class="aspect-square rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-600 flex items-center justify-center hover:border-indigo-500 transition-colors cursor-pointer group"
                @click.stop="triggerFileInput"
              >
                <div class="text-center">
                  <i class="fas fa-plus text-2xl text-gray-400 dark:text-gray-500 group-hover:text-indigo-500 transition-colors"></i>
                  <span class="block text-xs text-gray-500 dark:text-gray-400 mt-1">Add More</span>
                  <span class="block text-[10px] text-gray-400 dark:text-gray-500">{{ 5 - imagePreviews.length }} slots left</span>
                </div>
              </div>
            </div>

            <!-- Empty State -->
            <div v-else class="py-8">
              <div class="flex flex-col items-center justify-center">
                <i class="fas fa-cloud-upload-alt text-4xl text-gray-400 dark:text-gray-500 mb-3"></i>
                <p class="text-sm text-gray-600 dark:text-gray-400 mb-1 font-medium">
                  Drag & drop images here or click to browse
                </p>
                <p class="text-xs text-gray-500 dark:text-gray-400">
                  JPEG, PNG, WebP (max 10MB each) • Max 5 images
                </p>
                <button
                  type="button"
                  class="mt-4 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm rounded-lg transition-colors"
                  @click.stop="triggerFileInput"
                >
                  <i class="fas fa-upload mr-2"></i> Select Images
                </button>
              </div>
            </div>
          </div>

          <!-- Hidden file input -->
          <input
            ref="fileInput"
            type="file"
            accept="image/*"
            multiple
            @change="handleImageUpload"
            class="hidden"
          />

          <!-- Upload Progress Indicator -->
          <div v-if="isUploadingImages" class="mt-3">
            <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
              <i class="fas fa-spinner fa-spin text-indigo-500"></i>
              <span>Processing images...</span>
            </div>
          </div>
        </div>

        <VeeErrorMessage name="images" class="text-red-500 text-xs mt-1 block" />
      </div>

      <!-- Submit Buttons -->
      <div class="flex gap-4 pt-4 border-t border-gray-200 dark:border-gray-700">
        <button
          type="submit"
          :disabled="isSubmitting || isSubmittingLocal || isUploadingImages"
          class="flex-1 py-3 bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-400 text-white text-sm font-semibold rounded-lg shadow flex items-center justify-center gap-2 transition-colors"
        >
          <i v-if="isSubmitting || isSubmittingLocal || isUploadingImages" class="fas fa-spinner fa-spin"></i>
          <span v-else><i class="fas fa-plus"></i> Create Event</span>
        </button>
        <NuxtLink
          to="/"
          class="flex-1 py-3 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-200 text-sm font-semibold rounded-lg text-center transition-colors"
        >
          Cancel
        </NuxtLink>
      </div>
    </VeeForm>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from "vue";
import { useMutation } from "@vue/apollo-composable";
import { CREATE_EVENT_MUTATION } from "~/graphql/eventMutations";
import { useAuthStore } from "~/stores/auth";
import { useRouter } from "vue-router";
import * as yup from "yup";
import { useToast } from 'vue-toastification'

const toast = useToast();
const router = useRouter();
const authStore = useAuthStore();

// Form state
const form = reactive<{ [key: string]: any }>({
  title: "",
  description: "",
  category_id: "",
  price: 0,
  is_free: true,
  venue: "",
  address: "",
  latitude: 9.0319,
  longitude: 38.7468,
  event_date: "",
  start_time: "",
  end_time: "",
  status: "published",
  images: [], // Array of base64 images
});

const fileInput = ref<HTMLInputElement | null>(null);
const selectedImages = ref<File[]>([]);
const imagePreviews = ref<string[]>([]);
const tagInput = ref("");
const selectedTags = ref<string[]>([]);
const error = ref<string | null>(null);
const isSubmittingLocal = ref(false);
const isUploadingImages = ref(false);

// Validation schema
const schema = yup.object({
  title: yup
    .string()
    .required("Title is required")
    .min(3, "Title must be at least 3 characters")
    .max(255, "Title must be less than 255 characters"),
  description: yup
    .string()
    .required("Description is required")
    .min(20, "Description must be at least 20 characters"),
  category_id: yup.string().required("Category is required"),
  venue: yup.string().required("Venue is required"),
  address: yup.string().required("Address is required"),
  event_date: yup.string().required("Event date is required"),
  images: yup.array().min(1, "At least one image is required"),
});

// Image handling
const handleAreaClick = (event: MouseEvent) => {
  const target = event.target as HTMLElement;
  if (!target.closest('button') && imagePreviews.value.length < 5) {
    triggerFileInput();
  }
};

const triggerFileInput = () => {
  if (selectedImages.value.length < 5) {
    fileInput.value?.click();
  } else {
    toast.warning("Maximum 5 images allowed");
  }
};

const handleDragOver = (event: DragEvent) => {
  event.preventDefault();
  const target = event.currentTarget as HTMLElement;
  target.classList.add('border-indigo-500', 'bg-indigo-50', 'dark:bg-indigo-900/10');
};

const handleImageUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const files = target.files;
  if (files && files.length > 0) {
    const fileArray = Array.from(files);
    processImageFiles(fileArray);
  }
  target.value = "";
};

const handleDrop = (event: DragEvent) => {
  event.preventDefault();
  const target = event.currentTarget as HTMLElement;
  target.classList.remove('border-indigo-500', 'bg-indigo-50', 'dark:bg-indigo-900/10');
  
  const files = event.dataTransfer?.files;
  if (files && files.length > 0) {
    const fileArray = Array.from(files);
    processImageFiles(fileArray);
  }
};

const processImageFiles = (files: File[]) => {
  const currentCount = selectedImages.value.length;
  const maxSlots = 5 - currentCount;
  
  if (maxSlots <= 0) {
    toast.warning("Maximum 5 images already uploaded");
    return;
  }

  const validFiles = files
    .slice(0, maxSlots)
    .filter(validateImage);

  if (validFiles.length === 0) {
    if (files.length > 0) {
      toast.error("Please upload valid image files (JPEG, PNG, WebP, max 10MB each)");
    }
    return;
  }

  error.value = null;
  isUploadingImages.value = true;

  let processedCount = 0;
  validFiles.forEach((file) => {
    selectedImages.value.push(file);
    
    const reader = new FileReader();
    reader.onload = (e) => {
      imagePreviews.value.push(e.target?.result as string);
      processedCount++;
      
      if (processedCount === validFiles.length) {
        isUploadingImages.value = false;
        
        toast.success(`${validFiles.length} image${validFiles.length > 1 ? 's' : ''} added successfully!`, {
          timeout: 2000
        });
      }
    };
    reader.onerror = () => {
      processedCount++;
      if (processedCount === validFiles.length) {
        isUploadingImages.value = false;
      }
    };
    reader.readAsDataURL(file);
  });
};

const validateImage = (file: File): boolean => {
  if (file.size > 10 * 1024 * 1024) {
    toast.error(`"${file.name}" is larger than 10MB`);
    return false;
  }
  
  const validTypes = ["image/jpeg", "image/png", "image/webp"];
  if (!validTypes.includes(file.type)) {
    toast.error(`"${file.name}" must be JPEG, PNG, or WebP`);
    return false;
  }
  
  const isDuplicate = selectedImages.value.some(existing => 
    existing.name === file.name && existing.size === file.size
  );
  
  if (isDuplicate) {
    toast.warning(`"${file.name}" is already added`);
    return false;
  }
  
  return true;
};

const removeImage = (index: number) => {
  selectedImages.value.splice(index, 1);
  imagePreviews.value.splice(index, 1);
  toast.info("Image removed", { timeout: 1500 });
};

// Tag management
const addTag = () => {
  const normalized = tagInput.value.trim().toLowerCase();
  if (normalized && !selectedTags.value.includes(normalized) && selectedTags.value.length < 10) {
    selectedTags.value.push(normalized);
    tagInput.value = "";
  }
};

const removeTag = (index: number) => {
  selectedTags.value.splice(index, 1);
};

// Location update
const updateLocation = (location: {
  latitude: number;
  longitude: number;
  address: string;
}) => {
  form.latitude = location.latitude;
  form.longitude = location.longitude;
  form.address = location.address;
};

// Convert all images to base64
const imagesToBase64 = (files: File[]): Promise<string[]> => {
  return Promise.all(
    files.map((file) => {
      return new Promise<string>((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result as string);
        reader.onerror = (error) => reject(error);
      });
    })
  );
};

// Mutation
const { mutate: createEventMutation } = useMutation(CREATE_EVENT_MUTATION);

// Submit handler
const handleSubmit = async () => {
  try {
    error.value = null;

    // Validate images
    if (selectedImages.value.length === 0) {
      error.value = "Please upload at least one event image";
      toast.error("Please upload at least one event image");
      return;
    }

    isSubmittingLocal.value = true;
    isUploadingImages.value = true;

    // Convert ALL images to base64
    const base64Images = await imagesToBase64(selectedImages.value);

    // Prepare event data with all images
    const eventData = {
      title: form.title,
      description: form.description,
      category_id: form.category_id || null,
      price: form.is_free ? 0 : Number(form.price),
      is_free: form.is_free,
      venue: form.venue,
      address: form.address,
      latitude: form.latitude,
      longitude: form.longitude,
      event_date: form.event_date,
      start_time: form.start_time || null,
      end_time: form.end_time || null,
      status: form.status,
      images: base64Images, // All images - first one will be featured
      tags: selectedTags.value,
    };

    console.log("Creating event with all images:", eventData);

    // Get auth token
    let rawToken = authStore.token;
    if (!rawToken && typeof window !== 'undefined') {
      rawToken = localStorage.getItem("auth_token");
    }

    let authHeader = "";
    if (rawToken) {
      authHeader = rawToken.startsWith("Bearer ") ? rawToken : `Bearer ${rawToken}`;
    }

    // Single mutation with all images
    const result = await createEventMutation(
      { input: eventData },
      {
        context: {
          headers: {
            Authorization: authHeader,
          },
        },
      }
    );

    console.log("Mutation result:", result);

    if (result?.data?.createEvent?.id) {
      const eventId = result.data.createEvent.id;
      toast.success(`Event created successfully with ${base64Images.length} images! 🎉`);
      
      setTimeout(() => {
        router.push(`/events/${eventId}`);
      }, 2000);
    } else {
      const errorMsg = result?.data?.createEvent?.message || "Failed to create event. Please try again.";
      error.value = errorMsg;
      toast.error(errorMsg);
    }
  } catch (err: any) {
    error.value = err.message || "An unexpected error occurred";
    console.error("Submit error:", err);
    toast.error("Failed to create event");
  } finally {
    isSubmittingLocal.value = false;
    isUploadingImages.value = false;
  }
};

// Initialize date
onMounted(() => {
  form.event_date = new Date().toISOString().split("T")[0];
  authStore.loadAuth();
});
</script>

<style scoped>
/* Smooth transitions */
* {
  transition: background-color 0.15s, border-color 0.15s, opacity 0.15s;
}

/* Image hover effects */
.group:hover .group-hover\\:opacity-100 {
  opacity: 1;
}

/* Drag over effect */
.border-indigo-500 {
  border-color: rgb(99 102 241) !important;
}

.bg-indigo-50 {
  background-color: rgb(238 242 255) !important;
}

.dark\:bg-indigo-900\/10 {
  background-color: rgba(49, 46, 129, 0.1) !important;
}
</style>