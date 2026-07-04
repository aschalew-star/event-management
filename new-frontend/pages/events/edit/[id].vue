<template>
  <div class="container mx-auto max-w-3xl py-8 px-4">
    <div class="flex items-center justify-between mb-8">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        Edit Event
      </h1>
      <NuxtLink
        :to="`/events/${eventId}`"
        class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300 flex items-center gap-2"
      >
        <Icon name="lucide:eye" class="w-4 h-4" />
        View Event
      </NuxtLink>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center py-16">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="queryError || !event" class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl">
      <Icon name="lucide:alert-circle" class="w-20 h-20 text-red-400 mx-auto mb-4" />
      <h3 class="text-2xl font-semibold text-gray-700 dark:text-gray-200 mb-2">Failed to Load Event</h3>
      <p class="text-gray-500 dark:text-gray-400 mb-6">{{ queryError?.message || 'Event not found' }}</p>
      <NuxtLink
        to="/my-events"
        class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition-colors"
      >
        <Icon name="lucide:arrow-left" class="w-5 h-5" />
        Back to My Events
      </NuxtLink>
    </div>

    <!-- Edit Form -->
    <div v-else>
      <div
        v-if="isSubmittingLocal"
        class="mb-6 p-4 bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300 rounded-lg text-sm flex items-center gap-3 border border-blue-200 dark:border-blue-800"
      >
        <Icon name="lucide:loader-2" class="w-5 h-5 animate-spin" />
        <span>Updating event... Please wait.</span>
      </div>

      <div
        v-if="error"
        class="mb-6 p-4 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg text-sm border border-red-200 dark:border-red-800"
      >
        ⚠️ {{ error }}
      </div>

      <VeeForm
        v-slot="{ errors, isSubmitting }"
        @submit="handleSubmit"
        class="space-y-6"
        :validation-schema="schema"
      >
        <div>
          <label
            class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1"
          >
            Event Title *
          </label>
          <VeeField
            name="title"
            type="text"
            class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-4 py-2.5 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Enter event title"
            v-model="form.title"
          />
          <VeeErrorMessage name="title" class="text-red-500 text-xs mt-1 block" />
        </div>

        <div>
          <label
            class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1"
          >
            Description *
          </label>
          <VeeField
            name="description"
            as="textarea"
            rows="5"
            class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-4 py-2.5 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Describe your event..."
            v-model="form.description"
          />
          <VeeErrorMessage
            name="description"
            class="text-red-500 text-xs mt-1 block"
          />
        </div>

        <BaseCategorySelect v-model="form.category_id" />

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label
              class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2"
            >
              Price Rules
            </label>
            <div class="flex items-center gap-2 mb-2">
              <input
                id="isFreeCheck"
                type="checkbox"
                class="w-4 h-4 text-blue-600 rounded border-gray-300 focus:ring-blue-500"
                v-model="form.is_free"
              />
              <label
                for="isFreeCheck"
                class="text-sm text-gray-600 dark:text-gray-400"
              >
                Free Entry Event
              </label>
            </div>
            <VeeField
              v-if="!form.is_free"
              name="price"
              type="number"
              step="0.01"
              min="0"
              class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-4 py-2 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Price Amount"
              v-model="form.price"
            />
            <VeeErrorMessage
              name="price"
              class="text-red-500 text-xs mt-1 block"
            />
          </div>

          <div>
            <label
              class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1"
            >
              Venue Name *
            </label>
            <VeeField
              name="venue"
              type="text"
              class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-4 py-2.5 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Venue name"
              v-model="form.venue"
            />
            <VeeErrorMessage
              name="venue"
              class="text-red-500 text-xs mt-1 block"
            />
          </div>
        </div>

        <div
          class="p-4 bg-gray-50 dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-700"
        >
          <label
            class="block text-sm font-bold text-gray-800 dark:text-gray-200 mb-1"
          >
            Location *
          </label>
          <VeeField
            name="address"
            type="text"
            class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-4 py-2.5 text-sm text-gray-900 dark:text-white mb-3 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Full address"
            v-model="form.address"
          />
          <VeeErrorMessage
            name="address"
            class="text-red-500 text-xs mt-1 block"
          />

          <EventMapSingle
            :latitude="form.latitude"
            :longitude="form.longitude"
            mode="select"
            @update:location="updateLocation"
          />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label
              class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1"
            >
              Event Date *
            </label>
            <VeeField
              name="event_date"
              type="date"
              class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              v-model="form.event_date"
            />
            <VeeErrorMessage
              name="event_date"
              class="text-red-500 text-xs mt-1 block"
            />
          </div>
          <div>
            <label
              class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1"
            >
              Start Time
            </label>
            <VeeField
              name="start_time"
              type="time"
              class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              v-model="form.start_time"
            />
          </div>
          <div>
            <label
              class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1"
            >
              End Time
            </label>
            <VeeField
              name="end_time"
              type="time"
              class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              v-model="form.end_time"
            />
          </div>
        </div>

        <div>
          <label
            class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1"
          >
            Status
          </label>
          <div class="flex gap-4">
            <label class="flex items-center gap-2">
              <input
                type="radio"
                value="published"
                v-model="form.status"
                class="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
              />
              <span class="text-sm text-gray-700 dark:text-gray-300">Published</span>
            </label>
            <label class="flex items-center gap-2">
              <input
                type="radio"
                value="draft"
                v-model="form.status"
                class="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
              />
              <span class="text-sm text-gray-700 dark:text-gray-300">Draft</span>
            </label>
          </div>
        </div>

        <div>
          <label
            class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1"
          >
            Tags
          </label>
          <div
            class="bg-white dark:bg-gray-800 rounded-lg border border-gray-300 dark:border-gray-700 px-1 py-0.5 focus-within:ring-2 focus-within:ring-blue-500"
          >
            <input
              v-model="tagInput"
              @keydown.enter.prevent="addTag"
              placeholder="Type a tag and hit Enter..."
              class="w-full px-3 py-2 text-sm bg-transparent border-none outline-none text-gray-900 dark:text-white"
            />
          </div>
          <div class="flex flex-wrap gap-1.5 mt-2">
            <span
              v-for="(tag, i) in selectedTags"
              :key="i"
              class="inline-flex items-center gap-1 text-xs px-2.5 py-1 rounded-full bg-blue-50 dark:bg-blue-900/40 text-blue-600 dark:text-blue-400 font-medium"
            >
              {{ tag }}
              <button
                type="button"
                @click="removeTag(i)"
                class="hover:text-red-500 text-[10px]"
              >
                <Icon name="lucide:x" class="w-3 h-3" />
              </button>
            </span>
          </div>
        </div>

        <div>
          <label
            class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1"
          >
            Featured Image *
          </label>
          <div
            class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-6 text-center hover:border-blue-500 transition-colors"
            @dragover.prevent
            @drop.prevent="handleDrop"
          >
            <div v-if="imagePreview" class="mb-4">
              <div class="relative inline-block">
                <img
                  :src="imagePreview"
                  alt="Featured image preview"
                  class="max-h-48 rounded-lg object-cover"
                />
                <button
                  type="button"
                  @click="removeImage"
                  class="absolute -top-2 -right-2 p-1 bg-red-500 text-white rounded-full hover:bg-red-600 transition-colors"
                >
                  <Icon name="lucide:x" class="w-3 h-3" />
                </button>
              </div>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-2">
                Click the X to remove and upload a new image
              </p>
            </div>

            <div v-else>
              <input
                type="file"
                accept="image/*"
                @change="handleImageUpload"
                class="hidden"
                ref="fileInput"
              />
              <button
                type="button"
                @click="fileInput?.click()"
                class="px-4 py-2 bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200 text-sm rounded font-medium transition-colors"
              >
                <Icon name="lucide:cloud-upload" class="w-4 h-4 inline mr-2" />
                Choose Image
              </button>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-2">
                JPEG, PNG, WebP (max 10MB)
              </p>
            </div>
          </div>

          <VeeField
            name="featured_image"
            type="hidden"
            v-model="form.featured_image"
          />
          <VeeErrorMessage
            name="featured_image"
            class="text-red-500 text-xs mt-1 block"
          />
        </div>

        <div
          class="flex gap-4 pt-4 border-t border-gray-200 dark:border-gray-700"
        >
          <button
            type="submit"
            :disabled="isSubmitting || isSubmittingLocal"
            class="flex-1 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white text-sm font-semibold rounded-lg shadow flex items-center justify-center gap-2 transition-colors"
          >
            <Icon
              v-if="isSubmitting || isSubmittingLocal"
              name="lucide:loader-2"
              class="w-5 h-5 animate-spin"
            />
            <Icon v-else name="lucide:save" class="w-5 h-5" />
            <span v-if="isSubmitting || isSubmittingLocal">Updating...</span>
            <span v-else>Update Event</span>
          </button>
          <NuxtLink
            :to="`/my-events`"
            class="flex-1 py-3 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-200 text-sm font-semibold rounded-lg text-center transition-colors"
          >
            Cancel
          </NuxtLink>
        </div>
      </VeeForm>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useQuery, useMutation } from "@vue/apollo-composable";
import { useAuthStore } from "~/stores/auth";
import { useToast } from "vue-toastification";
import * as yup from "yup";
import { GET_EVENT_BY_ID } from "~/graphql/eventQueries";
import { UPDATE_EVENT_MUTATION } from "~/graphql/eventMutations";

definePageMeta({
  middleware: 'auth'
});

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

const eventId = computed(() => {
  const id = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id;
  return id || null;
});

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
  featured_image: undefined,
});

const fileInput = ref<HTMLInputElement | null>(null);
const selectedImage = ref<File | null>(null);
const imagePreview = ref<string | null>(null);
const tagInput = ref("");
const selectedTags = ref<string[]>([]);
const error = ref<string | null>(null);
const isSubmittingLocal = ref(false);
const originalFeaturedImage = ref<string | null>(null);

// Fetch event data
const { 
  result: eventResult, 
  loading, 
  error: queryError,
  refetch,
  onError 
} = useQuery(
  GET_EVENT_BY_ID,
  () => ({ id: eventId.value }),
  {
    fetchPolicy: 'network-only',
    skip: () => !eventId.value,
  }
);

onError((error) => {
  console.error('Error fetching event:', error);
  toast.error('Failed to load event data');
});

const event = computed(() => eventResult.value?.events_by_pk);

// Validation schema
const schema = yup.object({
  title: yup
    .string()
    .required("Title is required")
    .min(3, "Title must be at least 3 characters"),
  description: yup
    .string()
    .required("Description is required")
    .min(20, "Description must be at least 20 characters"),
  category_id: yup.string().required("Category is required"),
  venue: yup.string().required("Venue is required"),
  address: yup.string().required("Address is required"),
  event_date: yup.string().required("Event date is required"),
  featured_image: yup.mixed().required("Featured image is required"),
});

// Populate form with event data
const populateForm = () => {
  if (!event.value) return;

  form.title = event.value.title || "";
  form.description = event.value.description || "";
  form.category_id = event.value.category_id || "";
  form.is_free = event.value.is_free ?? true;
  form.price = event.value.price || 0;
  form.venue = event.value.venue || "";
  form.address = event.value.address || "";
  form.latitude = event.value.latitude || 9.0319;
  form.longitude = event.value.longitude || 38.7468;
  form.event_date = event.value.event_date ? event.value.event_date.split('T')[0] : "";
  form.start_time = event.value.start_time || "";
  form.end_time = event.value.end_time || "";
  form.status = event.value.status || "published";
  
  // Store original image for reference
  originalFeaturedImage.value = event.value.featured_image || null;
  
  // If there's a featured image, show it as preview
  if (event.value.featured_image) {
    imagePreview.value = event.value.featured_image;
  }

  // Parse tags if they exist
  if (event.value.tags) {
    try {
      const tags = typeof event.value.tags === 'string' 
        ? JSON.parse(event.value.tags) 
        : event.value.tags;
      selectedTags.value = Array.isArray(tags) ? tags : [];
    } catch {
      selectedTags.value = [];
    }
  }
};

// Watch for event data and populate form
watch(event, (newEvent) => {
  if (newEvent) {
    populateForm();
  }
}, { immediate: true });

// Image handlers
const handleImageUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (file) processImageFile(file);
  target.value = "";
};

const handleDrop = (event: DragEvent) => {
  const file = event.dataTransfer?.files?.[0];
  if (file) processImageFile(file);
};

const processImageFile = (file: File) => {
  if (!validateImage(file)) return;

  error.value = null;
  selectedImage.value = file;
  form.featured_image = file;

  const reader = new FileReader();
  reader.onload = (e) => {
    imagePreview.value = e.target?.result as string;
  };
  reader.readAsDataURL(file);
};

const validateImage = (file: File): boolean => {
  if (file.size > 10 * 1024 * 1024) {
    error.value = "Image size must be less than 10MB";
    return false;
  }
  const validTypes = ["image/jpeg", "image/png", "image/webp"];
  if (!validTypes.includes(file.type)) {
    error.value = "Only JPEG, PNG, and WebP images are allowed";
    return false;
  }
  return true;
};

const removeImage = () => {
  selectedImage.value = null;
  imagePreview.value = originalFeaturedImage.value;
  form.featured_image = originalFeaturedImage.value || undefined;
  if (fileInput.value) {
    fileInput.value.value = "";
  }
};

// Tag handlers
const addTag = () => {
  const normalized = tagInput.value.trim().toLowerCase();
  if (normalized && !selectedTags.value.includes(normalized)) {
    selectedTags.value.push(normalized);
  }
  tagInput.value = "";
};

const removeTag = (index: number) => {
  selectedTags.value.splice(index, 1);
};

// Location handler
const updateLocation = (location: {
  latitude: number;
  longitude: number;
  address: string;
}) => {
  form.latitude = location.latitude;
  form.longitude = location.longitude;
  form.address = location.address;
};

// Image to base64 conversion
const imageToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = (error) => reject(error);
  });
};

// Update mutation
const { mutate: updateEventMutation } = useMutation(UPDATE_EVENT_MUTATION);

const handleSubmit = async () => {
  try {
    error.value = null;
    
    // Validate required fields
    if (!form.title || !form.description || !form.venue || !form.address || !form.event_date) {
      error.value = "Please fill in all required fields";
      return;
    }

    isSubmittingLocal.value = true;

    let featuredImage = originalFeaturedImage.value;
    
    // If new image selected, convert to base64
    if (selectedImage.value) {
      featuredImage = await imageToBase64(selectedImage.value);
    }

    const input = {
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
      featured_image: featuredImage,
      tags: selectedTags.value,
    };

    console.log("Updating event with ID:", eventId.value);
    console.log("Update input:", input);

    const result = await updateEventMutation({
      id: eventId.value,
      input: input,
    });

    console.log("Update result:", result);

    // Check for success based on the mutation type
    if (result?.data?.update_events_by_pk?.id) {
      toast.success('Event updated successfully!');
      router.push(`/events/${result.data.update_events_by_pk.id}`);
    } else if (result?.data?.update_events?.returning?.length > 0) {
      toast.success('Event updated successfully!');
      router.push(`/events/${result.data.update_events.returning[0].id}`);
    } else {
      error.value = "Failed to update event. No data returned.";
      toast.error(error.value);
    }
  } catch (err: any) {
    error.value = err.message || "An unexpected error occurred";
    console.error("Submit error:", err);
    toast.error(error.value);
  } finally {
    isSubmittingLocal.value = false;
  }
};

// Refetch on mount
onMounted(async () => {
  if (eventId.value) {
    await refetch();
  }
});
</script>